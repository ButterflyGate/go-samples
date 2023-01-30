package main

import (
	"embed"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/xerrors"
)

//go:embed textfile.txt
var F embed.FS

func main() {
	fmt.Printf("%v", encrypto("here is our free space.", 10000))
	fmt.Printf("%v", dectypto([]int32{3997, 15031, 15426, 44437, 56963, 57278, 72064, 104902, 105910, 122056, 123942, 124632, 128637, 138160}, 10000))

}

func readtest() {
	fName, _ := os.Executable()
	cfName := fmt.Sprintf("%s.temporalycopy", fName)
	f, _ := os.Open(fName)
	cf, _ := os.Create(cfName)
	err := cf.Chmod(755)

	io.Copy(f, cf)

	s := "here is our free space"
	seek, err := find(f, s)
	if err != nil {
		fmt.Printf("%d: %+v\n", seek, err)
		return
	}
	fmt.Printf("found at %x\n", seek)
}

func find(f *os.File, s string) (int64, error) {
	letters := make([]byte, 20005)
	times := 10000

	seek := int64(0)
	for i := 0; i < times; i++ {
		c, err := f.ReadAt(letters, seek)
		if err != nil {
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

func encrypto(str string, rng int) []int32 {
	encrypted := make([]int32, 0, len(str))
	rand.Seed(int64(random()))

	bytes := []rune(str)

	cnt := int32(0)
	for i := 0; i < len(bytes) && cnt >= 0; cnt++ {
		r := rand.Int31n(int32(rng))
		k := bytes[i]
		_ = k

		if bytes[i] == r {
			encrypted = append(encrypted, cnt)
			i++
		}

	}

	return encrypted
}

func dectypto(num []int32, rng int) string {
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
