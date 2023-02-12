package main

import (
	"errors"
	"testing"

	//"fmt"
	"github.com/stretchr/testify/require"
	"github.com/udhos/equalfile"
	"os"
)

type testCase struct {
	caseName       string
	from           string
	to             string
	limit          int64
	offset         int64
	goldenFilePath string
	expectedError  error
}

func TestCopy(t *testing.T) {
	//PositiveCases
	inputPath := "./testdata/input.txt"
	outputPath := "out.txt"
	outputCopyPath := "testdata/input_copy.txt"
	cases := make([]testCase, 0)
	newCase := testCase{
		caseName:       "offset0_limit0",
		from:           inputPath,
		to:             outputPath,
		limit:          0,
		offset:         0,
		goldenFilePath: "testdata/out_offset0_limit0.txt",
		expectedError:  nil,
	}
	cases = append(cases, newCase)
	newCase = testCase{
		caseName:       "from==to",
		from:           inputPath,
		to:             inputPath,
		limit:          0,
		offset:         0,
		goldenFilePath: "testdata/out_offset0_limit0.txt",
		expectedError:  nil,
	}
	cases = append(cases, newCase)
	newCase = testCase{
		caseName:       "offset0_limit10",
		from:           inputPath,
		to:             outputPath,
		limit:          10,
		offset:         0,
		goldenFilePath: "testdata/my_out_offset0_limit10.txt",
		expectedError:  nil,
	}
	cases = append(cases, newCase)
	newCase = testCase{
		caseName:       "offset0_limit1000",
		from:           inputPath,
		to:             outputPath,
		limit:          1000,
		offset:         0,
		goldenFilePath: "testdata/my_out_offset0_limit1000.txt",
		expectedError:  nil,
	}
	cases = append(cases, newCase)
	newCase = testCase{
		caseName:       "offset0_limit10000",
		from:           inputPath,
		to:             outputPath,
		limit:          10000,
		offset:         0,
		goldenFilePath: "testdata/out_offset0_limit10000.txt",
		expectedError:  nil,
	}
	cases = append(cases, newCase)
	newCase = testCase{
		caseName:       "offset100_limit1000",
		from:           inputPath,
		to:             outputPath,
		limit:          1000,
		offset:         100,
		goldenFilePath: "testdata/my_out_offset100_limit1000.txt",
		expectedError:  nil,
	}
	cases = append(cases, newCase)
	newCase = testCase{
		caseName:       "offset6000_limit1000",
		from:           inputPath,
		to:             outputPath,
		limit:          1000,
		offset:         6000,
		goldenFilePath: "testdata/my_out_offset6000_limit1000.txt",
		expectedError:  nil,
	}
	cases = append(cases, newCase)

	var outputPathVariant string
	for _, tc := range cases {
		tc := tc
		t.Run(tc.caseName, func(t *testing.T) {
			//fmt.Println(tc)
			err := Copy(tc.from, tc.to, tc.offset, tc.limit)
			require.NoError(t, err)
			compare := equalfile.New(nil, equalfile.Options{})
			if tc.from != tc.to {
				outputPathVariant = tc.to
			} else {
				outputPathVariant = outputCopyPath
			}
			equal, err := compare.CompareFile(outputPathVariant, tc.goldenFilePath)
			if err != nil {
				require.Fail(t, "file comparsion error")
			}
			err = os.Remove(outputPathVariant)
			if err != nil {
				require.Fail(t, "test case output clear error")
			}
			require.True(t, equal)
		})
	}
	//Negative Cases
	cases = make([]testCase, 0)
	newCase = testCase{
		caseName:       "from pass is null",
		from:           "",
		to:             outputPath,
		limit:          0,
		offset:         0,
		goldenFilePath: "",
		expectedError:  ErrSourcePathIsNull,
	}
	cases = append(cases, newCase)
	newCase = testCase{
		caseName:       "to pass is null",
		from:           inputPath,
		to:             "",
		limit:          0,
		offset:         0,
		goldenFilePath: "",
		expectedError:  ErrTargetPathIsNull,
	}
	cases = append(cases, newCase)
	newCase = testCase{
		caseName:       "limit less 0",
		from:           inputPath,
		to:             outputPath,
		limit:          -5,
		offset:         0,
		goldenFilePath: "",
		expectedError:  ErrBadLimit,
	}
	cases = append(cases, newCase)
	newCase = testCase{
		caseName:       "offset less 0",
		from:           inputPath,
		to:             outputPath,
		limit:          0,
		offset:         -5,
		goldenFilePath: "",
		expectedError:  ErrBadOffset,
	}
	cases = append(cases, newCase)
	newCase = testCase{
		caseName:       "source is not exited",
		from:           "testdata/notExistedInput.txt",
		to:             outputPath,
		limit:          0,
		offset:         0,
		goldenFilePath: "",
		expectedError:  ErrSourceIsNotExisted,
	}
	cases = append(cases, newCase)
	newCase = testCase{
		caseName:       "unsupported file",
		from:           "testdata/UnsupportedInput.txt",
		to:             outputPath,
		limit:          0,
		offset:         0,
		goldenFilePath: "",
		expectedError:  ErrUnsupportedFile,
	}
	cases = append(cases, newCase)
	newCase = testCase{
		caseName:       "offset exceeds file size",
		from:           inputPath,
		to:             outputPath,
		limit:          0,
		offset:         1000000,
		goldenFilePath: "",
		expectedError:  ErrOffsetExceedsFileSize,
	}
	cases = append(cases, newCase)
	for _, tc := range cases {
		tc := tc
		t.Run(tc.caseName, func(t *testing.T) {
			err := Copy(tc.from, tc.to, tc.offset, tc.limit)
			require.Truef(t, errors.Is(err, tc.expectedError), "actual error %q", err)
		})
	}
}
