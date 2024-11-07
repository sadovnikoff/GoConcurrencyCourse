package filesystem

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

var now = time.Now

type Segment struct {
	file      *os.File
	directory string

	segmentSize    int
	maxSegmentSize int
}

func NewSegment(directory string, maxSegmentSize int) *Segment {
	return &Segment{
		directory:      directory,
		maxSegmentSize: maxSegmentSize,
	}
}

func (s *Segment) Write(data []byte) error {
	if s.file == nil || s.segmentSize >= s.maxSegmentSize {
		if err := s.createSegment(); err != nil {
			return fmt.Errorf("failed to create segment file: %w", err)
		}
	}

	writtenBytes, err := WriteFile(s.file, data)
	if err != nil {
		return fmt.Errorf("failed to write data to segment file: %w", err)
	}

	s.segmentSize += writtenBytes
	return nil
}

// TODO decompose this function

func (s *Segment) ReadAll() ([][]byte, error) {
	files, err := os.ReadDir(s.directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read WAL directory: %w", err)
	}

	re := regexp.MustCompile(`wal_(\d+)\.log`)
	fNames := map[int64]string{}
	orderedTimestamps := make([]int64, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			fName := file.Name()
			ts, err := s.extractTimestamp(re, fName)
			if err != nil {
				return nil, fmt.Errorf("failed to extract timestamp: %w", err)
			}

			orderedTimestamps = append(orderedTimestamps, ts)
			fNames[ts] = fName
		}
	}

	sort.Slice(orderedTimestamps, func(i, j int) bool {
		return orderedTimestamps[i] < orderedTimestamps[j]
	})

	rawData := make([][]byte, 0, len(orderedTimestamps))
	for _, ts := range orderedTimestamps {
		fullName := fmt.Sprintf("%s/%s", s.directory, fNames[ts])
		data, err := os.ReadFile(fullName)
		if err != nil {
			return nil, err
		}

		rawData = append(rawData, data)
	}

	return rawData, nil
}

func (s *Segment) createSegment() error {
	segmentName := fmt.Sprintf("%s/wal_%d.log", s.directory, now().UnixMilli())
	if s.file != nil {
		err := s.file.Close()
		if err != nil {
			return err
		}
	}

	file, err := CreateFile(segmentName)
	if err != nil {
		return err
	}

	s.file = file
	s.segmentSize = 0
	return nil
}

func (s *Segment) extractTimestamp(re *regexp.Regexp, filename string) (int64, error) {
	match := re.FindStringSubmatch(filename)
	if len(match) <= 1 {
		return 0, fmt.Errorf("failed to extract timestamp from file %s", filename)
	}

	timestamp, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return timestamp, nil
}
