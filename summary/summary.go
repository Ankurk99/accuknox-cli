// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Authors of KubeArmor

package summary

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	//opb "github.com/accuknox/auto-policy-discovery/src/protobuf/v1/observability"
	"github.com/fatih/color"
	"google.golang.org/grpc"
)

// Options Structure
type Options struct {
	GRPC      string
	Labels    string
	Namespace string
}

// StartSummary : Get summary on observability data
func StartSummary(o Options) error {
	gRPC := ""

	if o.GRPC != "" {
		gRPC = o.GRPC
	} else {
		if val, ok := os.LookupEnv("DISCOVERY_SERVICE"); ok {
			gRPC = val
		} else {
			gRPC = "localhost:9089"
		}
	}

	data := &LogsRequest{
		Label:     o.Labels,
		Namespace: o.Namespace,
	}

	// create a client
	conn, err := grpc.Dial(gRPC, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := NewSummaryClient(conn)

	//Fetch Summary Logs
	stream, err := client.FetchLogs(context.Background(), data)
	if err != nil {
		return errors.New("could not connect to the server. Possible troubleshooting:\n- Check if discovery engine is running\n- Create a portforward to discovery engine service using\n\t\033[1mkubectl port-forward -n explorer service/knoxautopolicy --address 0.0.0.0 --address :: 9089:9089\033[0m")
	}

	headerFmt := color.New(color.Underline).SprintfFunc()
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println("\n\n**********************************************************************")
		fmt.Println("\nPod Name : ", res.PodDetail)
		fmt.Println("\nNamespace : ", res.Namespace)
		//Print List of Processes
		fmt.Println("\nList of Processes (" + fmt.Sprint(len(res.ListOfProcess)) + ") :\n")
		tbl := Heading("SOURCE", "DESTINATION", "COUNT", "LAST UPDATED TIME", "STATUS")
		tbl.WithHeaderFormatter(headerFmt)
		for _, process := range res.ListOfProcess {
			for _, source := range process.ListOfDestination {
				tbl.AddRow(process.Source, source.Destination, source.Count, time.Unix(source.LastUpdatedTime, 0).Format("1-02-2006 15:04:05"), strings.ToUpper(source.Status))
			}
		}
		tbl.Print()

		//Print List of File System
		fmt.Println("\nList of File-system accesses (" + fmt.Sprint(len(res.ListOfFile)) + ") :\n")
		tbl = Heading("SOURCE", "DESTINATION", "COUNT", "LAST UPDATED TIME", "STATUS")
		tbl.WithHeaderFormatter(headerFmt)
		for _, file := range res.ListOfFile {
			for _, source := range file.ListOfDestination {
				tbl.AddRow(file.Source, source.Destination, source.Count, time.Unix(source.LastUpdatedTime, 0).Format("1-02-2006 15:04:05"), strings.ToUpper(source.Status))
			}
		}
		tbl.Print()

		//Print List of Network Connection
		fmt.Println("\nList of Network connections (" + fmt.Sprint(len(res.ListOfNetwork)) + ") :\n")
		tbl = Heading("SOURCE", "Protocol", "COUNT", "LAST UPDATED TIME", "STATUS")
		tbl.WithHeaderFormatter(headerFmt)
		for _, network := range res.ListOfNetwork {
			for _, source := range network.ListOfDestination {
				tbl.AddRow(network.Source, source.Destination, source.Count, time.Unix(source.LastUpdatedTime, 0).Format("1-02-2006 15:04:05"), strings.ToUpper(source.Status))
			}
		}
		tbl.Print()

		//Print Ingress Connections
		fmt.Printf("\nIngress Connections :\n\n")
		tbl = Heading("DESTINATION LABEL", "DESTINATION NAMESPACE", "PROTOCOL", "PORT", "COUNT", "LAST UPDATED TIME", "STATUS")
		tbl.WithHeaderFormatter(headerFmt)
		for _, ingress := range res.Ingress {
			tbl.AddRow(ingress.DestinationLabels, ingress.DestinationNamespace, ingress.Protocol, ingress.Port, ingress.Count, time.Unix(ingress.LastUpdatedTime, 0).Format("1-02-2006 15:04:05"), ingress.Status)
		}
		tbl.Print()

		//Print Egress Connections
		fmt.Printf("\nEgress Connections : \n\n")
		tbl = Heading("DESTINATION LABEL", "DESTINATION NAMESPACE", "PROTOCOL", "PORT", "COUNT", "LAST UPDATED TIME", "STATUS")
		tbl.WithHeaderFormatter(headerFmt)
		for _, egress := range res.Egress {
			tbl.AddRow(egress.DestinationLabels, egress.DestinationNamespace, egress.Protocol, egress.Port, egress.Count, time.Unix(egress.LastUpdatedTime, 0).Format("1-02-2006 15:04:05"), egress.Status)
		}
		tbl.Print()

		//Print System Incoming connections
		fmt.Println("\nList of Incoming connections (" + fmt.Sprint(len(res.InServerConn)) + ") :\n")
		tbl = Heading("PROTOCOL", "IP/PATH", "PORT", "LABELS", "NAMESPACE")
		tbl.WithHeaderFormatter(headerFmt)
		for _, inConn := range res.InServerConn {
			tbl.AddRow(inConn.Protocol, inConn.PodSvcIP, inConn.ServerPort, inConn.Labels, inConn.Namespace)
		}
		tbl.Print()

		//Print System Outgoing connections
		fmt.Println("\nList of Outgoing connections (" + fmt.Sprint(len(res.OutServerConn)) + ") :\n")
		tbl = Heading("PROTOCOL", "IP/PATH", "PORT", "LABELS", "NAMESPACE")
		tbl.WithHeaderFormatter(headerFmt)
		for _, outConn := range res.OutServerConn {
			tbl.AddRow(outConn.Protocol, outConn.PodSvcIP, outConn.ServerPort, outConn.Labels, outConn.Namespace)
		}
		tbl.Print()

	}
	return nil
}
