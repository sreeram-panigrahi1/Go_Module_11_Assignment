package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/sreeram-panigrahi1/Go_Module_11_Assignment/calculatorpb"

	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Start sending requests")

	cc, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	//Unary API function - Calculator
	Sum(c)

	// Server-Side Streaming function - PrimeNumbers
	PrimeNumbers(c)

	//Client-Side Streaming function - ComputeAverage
	ComputeAverage(c)

	//bidirectional streaming function - FindMaxNumber
	FindMaxNumber(c)

	fmt.Println("FINISHED")

}

func Sum(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("UNARY GRPC....")

	req := calculatorpb.SumRequest{
		Num1: 1.1,
		Num2: 2.2,
	}

	resp, err := c.Sum(context.Background(), &req)
	if err != nil {
		log.Fatalf("ERROR %v", err)
	}

	log.Printf("UNARY : Sum : %v", resp.GetSum())

}

func PrimeNumbers(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Serverside GRPC streaming")

	req := calculatorpb.PrimeNumbersRequest{
		Limit: 15,
	}

	respStream, err := c.PrimeNumbers(context.Background(), &req)
	if err != nil {
		log.Fatalf("ERROR %v", err)
	}

	for {
		msg, err := respStream.Recv()
		if err == io.EOF {
			//we have reached to the end of the file
			break
		}

		if err != nil {
			log.Fatalf("ERROR : %v", err)
		}

		log.Println("Server: Prime Number : ", msg.GetPrimeNum())
	}
}

func ComputeAverage(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("Client Side Streaming :")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("ERROR : %v", err)
	}

	requests := []*calculatorpb.ComputeAverageRequest{
		{
			Num: 2,
		},
		{
			Num: 2,
		},
		{
			Num: 2,
		},
		{
			Num: 2,
		},
	}

	for _, req := range requests {
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	log.Println("Average: ", resp.GetAvg())
}

func FindMaxNumber(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Bi Directional")

	requests := []*calculatorpb.FindMaxNumberRequest{
		{
			Num: 2,
		},
		{
			Num: 2,
		},
		{
			Num: 2,
		},
		{
			Num: 2,
		},
		{
			Num: 2,
		},
	}

	stream, err := c.FindMaxNumber(context.Background())
	if err != nil {
		log.Fatalf("ERROR : %v", err)
	}

	//wait channel to block receiver
	waitchan := make(chan struct{})

	go func(requests []*calculatorpb.FindMaxNumberRequest) {
		for _, req := range requests {

			err := stream.Send(req)
			if err != nil {
				log.Fatalf("ERROR : %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}(requests)

	go func() {
		for {

			resp, err := stream.Recv()
			if err == io.EOF {
				close(waitchan)
				return
			}

			if err != nil {
				log.Fatalf("ERROR : %v", err)
			}

			log.Printf("Max : %v\n", resp.GetMax())
		}
	}()

	//block until everything is finished
	<-waitchan
}
