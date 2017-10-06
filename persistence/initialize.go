package persistence

import (
	"log"
	"time"

	"github.com/vice-registry/vice-util/models"
	gocb "gopkg.in/couchbase/gocb.v1"
)

// InitViceCouchbase initializes an empty couchbase instance (e.g. creates admin account)
func InitViceCouchbase() {
	cluster, err := gocb.Connect("couchbase://" + couchbaseCredentials.Location)
	if err != nil {
		log.Fatalf("Cannot connect to couchbase: %s", err)
	}

	clusterManager := cluster.Manager(couchbaseCredentials.Username, couchbaseCredentials.Password)
	createBucket(cluster, clusterManager, "vice-users")
	createBucket(cluster, clusterManager, "vice-environments")
	createBucket(cluster, clusterManager, "vice-images")
	createBucket(cluster, clusterManager, "vice-deployments")

	createAdminUser(cluster)
}

func bucketExists(clusterManager *gocb.ClusterManager, bucketname string) bool {
	buckets, err := clusterManager.GetBuckets()
	if err != nil {
		log.Fatalf("Cannot get list of couchbase buckets: %s", err)
	}
	for i := range buckets {
		//log.Printf("bucket: %+v\n", buckets[i])
		if buckets[i].Name == bucketname {
			return true
		}
	}
	return false
}

func createBucket(cluster *gocb.Cluster, clusterManager *gocb.ClusterManager, bucketname string) {
	log.Printf("Create (if not exist) Couchbase bucket %s ...", bucketname)

	// check if bucket exists
	if bucketExists(clusterManager, bucketname) {
		log.Printf("Couchbase Bucket %s found, will not create it.", bucketname)
		return
	}

	// create bucket
	settings := gocb.BucketSettings{
		Name:     bucketname,
		Quota:    256,
		Replicas: 1,
		Type:     gocb.Couchbase,
		Password: couchbaseCredentials.Password,
	}
	err := clusterManager.InsertBucket(&settings)
	if err != nil {
		log.Fatalf("Cannot create bucket %s: %s", bucketname, err)
	}
	// wait for bucket
	if !bucketExists(clusterManager, bucketname) {
		log.Fatalf("Cannot find created bucket %s.", bucketname)
	}
	for i := 0; i < 10; i++ {
		testBucket, err := cluster.OpenBucket(bucketname, couchbaseCredentials.Password)
		if err != nil {
			log.Printf("Testing availability of bucket %s (try %d/10): %s", bucketname, i, err)
			time.Sleep(1000 * time.Millisecond)
		} else {
			testBucket.Close()
			// bucket available. go on.
			break
		}
	}

	// create bucket primary index
	bucket, err := cluster.OpenBucket(bucketname, couchbaseCredentials.Password)
	if err != nil {
		log.Fatalf("Cannot open bucket to create primary %s: %s", bucketname, err)
	}
	defer bucket.Close()
	bucketManager := bucket.Manager("", "")
	err = bucketManager.CreatePrimaryIndex("", true, false)
	if err != nil {
		log.Fatalf("Cannot create primary on %s: %s", bucketname, err)
	}
}

func createAdminUser(cluster *gocb.Cluster) {
	bucket, err := cluster.OpenBucket("vice-users", couchbaseCredentials.Password)
	if err != nil {
		log.Fatalf("Cannot create admin user, failed to open couchbase bucket vice-users: %s", err)
	}

	// define default admin user
	var admin models.User
	admin.ID = "admin"
	admin.Username = "admin"
	admin.Password = "admin"

	// try to insert admin user
	_, err = bucket.Insert(admin.ID, admin, 0)
	if err != nil {
		log.Printf("Admin user not created: %s", err)
	}
}
