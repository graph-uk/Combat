package LauncherServer

import (
	//"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (t *LauncherServer) launchWorkersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("LWHandler")
	LauncherServerKey := r.Header.Get("Authorization")
	log.Println("Auth=" + LauncherServerKey)
	if LauncherServerKey == t.config.LauncherSecretKey {
		path := strings.Split(r.URL.Path, "/")
		count, err := strconv.Atoi(path[len(path)-1])
		log.Println("count=" + strconv.Itoa(count))
		if err != nil {
			w.Write([]byte(`Incorrect count of workers: it should be int. For example "/LaunchWorkers/3" \r\n`))
			return
		}

		t.LaunchWorkers(count)
		log.Println(strconv.Itoa(count) + " workers launched")
	} else {
		log.Println(r.RemoteAddr + " Key is not valid.")
	}
}

func (t *LauncherServer) LaunchWorkers(CountOfWorkers int) error {
	os.Setenv("AWS_ACCESS_KEY_ID", t.config.AWS_ACCESS_KEY_ID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", t.config.AWS_SECRET_ACCESS_KEY)

	svc := ec2.New(session.New(&aws.Config{Region: aws.String(t.config.Region)}))

	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId: aws.String(t.config.AMIId),
		InstanceInitiatedShutdownBehavior: aws.String(t.config.ShutdownBehavior), //stop; terminate
		InstanceType:                      aws.String(t.config.InstanceType),
		SecurityGroupIds:                  aws.StringSlice([]string{t.config.SecurityGroupId}),
		KeyName:                           aws.String(t.config.KeyName),
		MinCount:                          aws.Int64(int64(CountOfWorkers)),
		MaxCount:                          aws.Int64(int64(CountOfWorkers)),
	})

	if err != nil {
		log.Println("Could not create instance", err)
		return err
	}

	// set name of instances.
	for _, curInstance := range runResult.Instances {
		_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
			Resources: []*string{curInstance.InstanceId},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(t.config.TagName),
				},
			},
		})
		if errtag != nil {
			log.Println("Could not create tags for instance", runResult.Instances[0].InstanceId, errtag)
			return errtag
		}
	}
	return nil
}
