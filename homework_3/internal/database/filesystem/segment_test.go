package filesystem

import (
	"os"
	"testing"
	"time"
)

func TestSegmentWrite(t *testing.T) {
	const testWALDirectory = "temp_test_data"
	err := os.Mkdir(testWALDirectory, os.ModePerm)
	if err != nil {
		t.Errorf("cannot create temporary dir for test data [%s]: %s", testWALDirectory, err)
	}

	defer func() {
		err := os.RemoveAll(testWALDirectory)
		if err != nil {
			t.Errorf("cannot remove temporary dir for test data [%s]: %s", testWALDirectory, err)
		}
	}()

	const maxSegmentSize = 10
	segment := NewSegment(testWALDirectory, maxSegmentSize)

	now = func() time.Time {
		return time.Unix(1, 0)
	}

	err = segment.Write([]byte("01234"))
	if err != nil {
		t.Errorf("cannot write test data: %s", err)
	}

	if segment.segmentSize != 5 {
		t.Errorf("wrong segment size: wannt 5, got %d", segment.segmentSize)
	}

	err = segment.Write([]byte("56789"))
	if err != nil {
		t.Errorf("cannot write test data: %s", err)
	}

	now = func() time.Time {
		return time.Unix(2, 0)
	}

	err = segment.Write([]byte("aaaaa"))
	if err != nil {
		t.Errorf("cannot write test data: %s", err)
	}

	if segment.segmentSize != 5 {
		t.Errorf("wrong segment size: wannt 5, got %d", segment.segmentSize)
	}

	stat, err := os.Stat(testWALDirectory + "/wal_1000.log")
	if err != nil {
		t.Errorf("cannot get file info [%s]", testWALDirectory+"/wal_1000.log")
	}

	if stat.Size() != 10 {
		t.Errorf("wrong file size: wannt 10, got %d", stat.Size())
	}

	stat, err = os.Stat(testWALDirectory + "/wal_2000.log")
	if err != nil {
		t.Errorf("cannot get file info [%s]", testWALDirectory+"/wal_2000.log")
	}

	if stat.Size() != 5 {
		t.Errorf("wrong file size: wannt 5, got %d", stat.Size())
	}

	err = segment.file.Close()
	if err != nil {
		t.Errorf("cannot close segment file [%s]", testWALDirectory+"/wal_2000.log")
	}
}

func TestReadAll(t *testing.T) {
	expectedSegmentsCount := 3
	testDir := "test_data"
	segment := NewSegment(testDir, 10)

	data, err := segment.ReadAll()
	if err != nil {
		t.Errorf("cannot read segmets data from [%s]", testDir)
	}

	if len(data) != expectedSegmentsCount {
		t.Errorf("wrong number of segments: wannt %d, got %d", expectedSegmentsCount, len(data))
	}

	if string(data[0]) != "0123" {
		t.Errorf("wrong segment data: wannt %s, got %s", "0123", string(data[0]))
	}

	if string(data[1]) != "4567" {
		t.Errorf("wrong segment data: wannt %s, got %s", "0123", string(data[0]))
	}

	if string(data[2]) != "89" {
		t.Errorf("wrong segment data: wannt %s, got %s", "0123", string(data[0]))
	}
}
