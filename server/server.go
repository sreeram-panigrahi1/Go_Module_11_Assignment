package main

import (
	"context"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/sreeram-panigrahi1/Go_Module_11_Assignment/calculatorpb"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (resp *calculatorpb.SumResponse, err error) {

	num1 := req.GetNum1()
	num2 := req.GetNum2()

	resp = &calculatorpb.SumResponse{
		Sum: num1 + num2,
	}
	return resp, nil
}

func (*server) PrimeNumbers(req *calculatorpb.PrimeNumbersRequest, resp calculatorpb.CalculatorService_PrimeNumbersServer) error {

	isPrime := func(num int64) bool {
		if num <= 1 {
			return false
		}
		limit := int64(math.Sqrt(float64(num)))
		for i := int64(2); i <= limit; i++ {
			if num%i == 0 {
				return false
			}
		}
		return true
	}

	limit := req.GetLimit()

	for i := int64(0); i <= limit; i++ {
		if isPrime(i) {
			res := calculatorpb.PrimeNumbersResponse{
				PrimeNum: i,
			}
			time.Sleep(1000 * time.Millisecond)
			resp.Send(&res)
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {

	var avg int64 = 0
	var count int64 = 0

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			//we have finished reading client stream
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Avg: avg / count,
			})
		}

		if err != nil {
			log.Fatalf("ERROR : %v", err)
		}

		num := msg.GetNum()
		count++
		avg += num
	}
}

func (*server) FindMaxNumber(stream calculatorpb.CalculatorService_FindMaxNumberServer) error {

	var max int64 = math.MinInt64

	for {

		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("ERROR : %v", err)
			return err
		}

		num := req.GetNum()

		if num > max {
			max = num
			sendErr := stream.Send(&calculatorpb.FindMaxNumberResponse{
				Max: max,
			})

			if sendErr != nil {
				log.Fatalf("error while sending response to Calculator Client : %v", err)
				return err
			}
		}
	}
	return nil
}

func main() {

	listen, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	log.Println("Starting Server")
	if err = s.Serve(listen); err != nil {
		log.Fatalf("ERROR : %v", err)
	}
}
