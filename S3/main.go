package main

func main() {
	svc := createSession()

	createBucket(svc, "se377krileyprogrammaticaccess", "us-east-1")

	listBucket(svc)

	createObject(svc, "se377krileyprogrammaticaccess", "test.txt")

	removeObject(svc, "se377krileyprogrammaticaccess", "test.txt")
}
