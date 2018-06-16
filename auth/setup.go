package users

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/UnnoTed/authenticaTed/errors"
	. "github.com/UnnoTed/authenticaTed/logger"

	log "github.com/Sirupsen/logrus"
	"github.com/c2h5oh/hide"
	"github.com/dustin/go-humanize"
	"golang.org/x/crypto/bcrypt"
)

const (
	minimumRounds = 14

	// seconds * (1s in nanosecond)
	ns            = 1000000000
	maxTime int64 = 3 * ns
)

type result struct {
	rounds int
	time   int64
}

var (
	p  = []byte("123456")
	cr = minimumRounds

	primes = []uint64{
		1199456261,
		8432571981118615261,
		1803515671,
		6604711928022367543,
		4365681490325845181,
	}
)

func benchmarkRoundTime(b *testing.B) {
	var err error

	b.ResetTimer()
	for i := 0; i < 1; i++ {
		_, err = bcrypt.GenerateFromPassword(p, cr)
	}
	b.StopTimer()

	if err != nil {
		panic(err)
	}
}

func Setup() *errors.Error {
	Logger.Info("[PASSWORD HASH ROUNDS]: Starting...")

	var (
		results []*result
		b       testing.BenchmarkResult
		t       int64
	)

	for {
		b = testing.Benchmark(benchmarkRoundTime)
		t = b.NsPerOp()

		if t >= maxTime {
			break
		} else {
			Logger.WithFields(log.Fields{
				"Rounds": cr,
				"Time":   humanize.Comma(t),
			}).Info("[PASSWORD HASH ROUNDS]: Not enough to ", maxTime/ns, " seconds")

			r := &result{
				rounds: cr,
				time:   t,
			}
			results = append(results, r)
			cr++
		}
	}

	Logger.WithFields(log.Fields{
		"Rounds": cr,
		"Time":   t / ns,
	}).Info("[PASSWORD HASH ROUNDS]: Done, round choosen is ", cr-1)
	log.Println(fmt.Sprintf("ROUND[%v] TIME[%v]", cr, t))

	var err error
	handleErr := func(f func(*big.Int) error, i *big.Int) {
		if err == nil {
			err = f(i)
		}

		if err != nil {
			log.Println("CAUSED BY", i)
		}
	}

	handleErr(hide.Default.SetInt32, new(big.Int).SetInt64(int64(primes[0])))
	handleErr(hide.Default.SetInt64, new(big.Int).SetInt64(int64(primes[1])))
	handleErr(hide.Default.SetUint32, new(big.Int).SetUint64(primes[2]))
	handleErr(hide.Default.SetUint64, new(big.Int).SetUint64(primes[3]))
	handleErr(hide.Default.SetXor, new(big.Int).SetUint64(primes[4]))

	/*d := gomail.NewDialer("smtp.example.com", 587, "user", "123456")
	_, err := d.Dial()
	if err != nil {
		panic(err)
	}*/

	return errors.FromErr(err)
}
