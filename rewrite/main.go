package main

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ButterflyGate/logger"
	"golang.org/x/xerrors"
)

//go:embed textfile.txt
var F embed.FS

func main() {

	firstStringMatch := decrypto([]int32{41608, 52728, 53406, 55231, 56963, 72064, 82762, 122056, 123121, 128312, 138943, 142020, 148154, 160243, 167466, 171191, 171576, 201908, 218988, 220626, 221161, 236283, 240908}, 10000)
	anytimeStringMatch := decrypto([]int32{3997, 15031, 15426, 44437, 56963, 57278, 72064, 104902, 105910, 122056, 123942, 124632, 128637, 138160}, 10000)

	lo := logger.DefaultOutputOption().HideTimestamp().HideLevel().HideLevel()
	lf := logger.DefaultFormatOption()
	l := logger.NewLogger(5, lo, lf)

	msg, err := F.ReadFile("textfile.txt")
	if err != nil {
		l.Error(err)
		return
	}
	l.Info(string(msg))

	fName, _ := os.Executable()
	cfName := fmt.Sprintf("%s.temporalycopy", fName)
	f, err := os.Open(fName)
	if err != nil {
		l.Error(err)
		return
	}
	defer f.Close()
	cf, err := os.Create(cfName)
	if err != nil {
		l.Error(err)
		return
	}
	defer cf.Close()
	_, err = io.Copy(cf, f)
	if err != nil {
		l.Error(err)
		return
	}
	err = cf.Chmod(0755)
	if err != nil {
		l.Error(err)
		return
	}

	execCnt := 0
	seek, err := find(cf, anytimeStringMatch)
	if seek == -1 {
		seek, err = find(cf, firstStringMatch)
		if err != nil {
			l.Error(err)
			return
		}
	} else {
		execCnt, err = getExecCnt(anytimeStringMatch, msg)
		l.Debug(execCnt, "execCnt")
	}
	if err != nil {
		l.Error(err)
		return
	}

	execCnt++
	fmt.Printf("これは %d 回目 の起動です。\n", execCnt)

	l.Debug(execCnt, "execCnt")
	rewriteMsg := fmt.Sprintf("%s%d", anytimeStringMatch, execCnt)
	err = rewrite(cf, firstStringMatch, rewriteMsg, seek)
	if err != nil {
		l.Error(err)
		err = os.Remove(cfName)
		if err != nil {
			l.Error(err)
		}
		return
	}

	err = os.Remove(fName)
	if err != nil {
		l.Error(err)
		return
	}
	err = os.Rename(cfName, fName)
	if err != nil {
		l.Error(err)
		return
	}
}

func find(f *os.File, s string) (int64, error) {
	letters := make([]byte, 20005)
	times := 10000

	seek := int64(0)
	for i := 0; i < times; i++ {
		c, err := f.ReadAt(letters, seek)
		if err != nil {
			if errors.As(io.EOF, &err) {
				break
			}
			return seek, err
		}
		if strings.Contains(string(letters), s) {
			idx := strings.LastIndex(string(letters), s)
			seek += int64(idx)
			return seek, nil
		}
		seek += int64(c)
	}
	return -1, xerrors.Errorf("could not find")
}

func rewrite(f *os.File, baseMatch, s string, seek int64) error {
	if len(s) > len(baseMatch) {
		return xerrors.Errorf("too long lettteral")
	}

	b := make([]byte, len(baseMatch))

	for i, v := range s {
		b[i] = byte(v)
	}
	size, err := f.WriteAt(b, seek)
	if err != nil {
		return err
	}
	if size != len(baseMatch) {
		return xerrors.Errorf("mismatch length")
	}
	return nil
}

func getExecCnt(baseMatch string, s []byte) (int, error) {
	s = s[len(baseMatch)-1:]
	nums := make([]byte, 0, 5)
	for _, v := range s {
		_, err := strconv.Atoi(string(v))
		if err == nil {
			nums = append(nums, v)
		}
	}
	n, err := strconv.Atoi(string(nums))
	if err != nil {
		return -1, xerrors.Errorf("failed convert to number: %w", err)
	}
	return n, nil
}

func decrypto(num []int32, rng int) string {
	decrypted := make([]rune, 0, len(num))
	rand.Seed(int64(random()))
	cnt := 0

	for i := 0; i < len(num) && cnt >= 0; cnt++ {
		r := rand.Intn(rng)
		if num[i] == int32(cnt) {
			decrypted = append(decrypted, rune(r))
			i++
		}
	}

	return string(decrypted)
}

// return 1
func random() int {
	rand.Seed(time.Hour.Nanoseconds())
	r := rand.Int()
	for i := 0; i < 100000; i++ {
		if r&1 == 0 {
			r = r / 2
		} else {
			r = r*3 + 1
		}
	}
	return strings.Count(fmt.Sprintf("%b", r), "1")
}
