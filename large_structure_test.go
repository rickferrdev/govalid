package govalid_test

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/rickferrdev/govalid"
)

type largeProfile struct {
	Email string
	Age   uint64
}

type largePayload struct {
	ID      uint64
	Name    string
	Profile largeProfile
	Scores  []int64
	Labels  map[string]string
	Matrix  [][]int
	Data    []byte
}

var largePayloadFields = []govalid.FieldSpec{
	govalid.Field("ID", govalid.NotZero()),
	govalid.Field("Name", govalid.Required(), govalid.StringMinLength(3)),
	govalid.Field("Profile.Email", govalid.Required(), govalid.StringEmail()),
	govalid.Field("Profile.Age", govalid.IntBetween(uint64(18), uint64(130))),
	govalid.Field("Scores", govalid.Required(), govalid.CollectionEach(govalid.IntBetween(int64(0), int64(100)))),
	govalid.Field("Labels", govalid.Required(), govalid.MapKeys(govalid.StringLowercase()), govalid.MapValues(govalid.StringRequired())),
	govalid.Field("Matrix", govalid.Required(), govalid.CollectionEach(govalid.CollectionEach(govalid.IntBetween(0, 1_000_000)))),
	govalid.Field("Data", govalid.Required(), govalid.BytesMinLength(1024), govalid.BytesUTF8(), govalid.BytesJSON()),
}

func buildLargePayload(scoreCount, labelCount, matrixRows, matrixColumns, dataSize int) largePayload {
	scores := make([]int64, scoreCount)
	for index := range scores {
		scores[index] = int64(index % 101)
	}

	labels := make(map[string]string, labelCount)
	for index := 0; index < labelCount; index++ {
		labels[fmt.Sprintf("key-%06d", index)] = fmt.Sprintf("value-%06d", index)
	}

	matrix := make([][]int, matrixRows)
	for row := range matrix {
		matrix[row] = make([]int, matrixColumns)
		for column := range matrix[row] {
			matrix[row][column] = row*matrixColumns + column
		}
	}

	minimumJSONSize := len(`{"data":""}`)
	if dataSize < minimumJSONSize {
		dataSize = minimumJSONSize
	}
	data := []byte(`{"data":"` + strings.Repeat("a", dataSize-minimumJSONSize) + `"}`)

	return largePayload{
		ID:      1,
		Name:    "large-payload",
		Profile: largeProfile{Email: "load-test@example.com", Age: 42},
		Scores:  scores,
		Labels:  labels,
		Matrix:  matrix,
		Data:    data,
	}
}

func TestLargeStructureValidation(t *testing.T) {
	payload := buildLargePayload(100_000, 10_000, 128, 128, 1<<20)

	if err := govalid.New().Validate(payload, largePayloadFields...); err != nil {
		t.Fatalf("large valid payload was rejected: %v", err)
	}
}

func TestLargeStructureCollectsIndependentIssues(t *testing.T) {
	payload := buildLargePayload(10_000, 1_000, 32, 32, 64<<10)
	payload.ID = 0
	payload.Profile.Email = "invalid-email"
	payload.Scores[0] = -1
	payload.Labels["UPPERCASE"] = ""
	payload.Data = []byte("not-json")

	err := govalid.New().Validate(payload, largePayloadFields...)
	var validationErr *govalid.FieldIssueError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected FieldIssueError, received %T", err)
	}
	if len(validationErr.RulesIssues) != 7 {
		t.Fatalf("expected 7 independent issues, received %d", len(validationErr.RulesIssues))
	}

	paths := make(map[string]int)
	for _, issue := range validationErr.RulesIssues {
		paths[issue.Path]++
	}
	for _, path := range []string{"ID", "Profile.Email", "Scores", "Labels", "Data"} {
		if paths[path] == 0 {
			t.Errorf("expected an issue for path %q, received %#v", path, paths)
		}
	}
	if formatted := err.Error(); len(formatted) > 16<<10 {
		t.Fatalf("large validation error should be summarized, received %d bytes", len(formatted))
	}
}

func TestLargeStructureConcurrentValidation(t *testing.T) {
	payload := buildLargePayload(25_000, 2_000, 64, 64, 256<<10)
	validator := govalid.New()

	workers := runtime.GOMAXPROCS(0) * 2
	if workers > 16 {
		workers = 16
	}
	if workers < 4 {
		workers = 4
	}

	start := make(chan struct{})
	failures := make(chan error, workers)
	var waitGroup sync.WaitGroup

	for worker := 0; worker < workers; worker++ {
		waitGroup.Add(1)
		go func(worker int) {
			defer waitGroup.Done()
			<-start
			if err := validator.Validate(payload, largePayloadFields...); err != nil {
				failures <- fmt.Errorf("worker %d: %w", worker, err)
			}
		}(worker)
	}

	close(start)
	waitGroup.Wait()
	close(failures)

	for err := range failures {
		t.Error(err)
	}
}

func BenchmarkLargeStructureValidation(b *testing.B) {
	payload := buildLargePayload(100_000, 10_000, 128, 128, 1<<20)
	validator := govalid.New()

	b.ReportAllocs()
	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		if err := validator.Validate(payload, largePayloadFields...); err != nil {
			b.Fatal(err)
		}
	}
}
