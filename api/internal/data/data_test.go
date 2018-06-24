package data

import (
	"testing"
	"github.com/gambarini/cabapi/tstutils"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestDb_FindTrips(t *testing.T) {

	err := tstutils.StartMySQL()
	defer tstutils.StopMySQL()

	if err != nil {
		t.Errorf("Failed to start MySQL, %s", err)
		return
	}

	err = tstutils.StartRedis()
	defer tstutils.StopRedis()

	if err != nil {
		t.Errorf("Failed to start Redis, %s", err)
		return
	}

	time.Sleep(time.Second * 10)

	err = tstutils.CreateDatabase()

	if err != nil {
		t.Errorf("Failed to create database, %s", err)
		return
	}

	db, err := NewDb()

	if err != nil {
		t.Errorf("Failed to create new db, %s", err)
		return
	}

	defer db.Close()

	_, err = db.db.Exec("CREATE TABLE `cab_trip_data` (" +
		"`medallion` text," +
		"`hack_license` text," +
		"`vendor_id` text," +
		"`rate_code` int(11) DEFAULT NULL," +
		"`store_and_fwd_flag` text," +
		"`pickup_datetime` datetime DEFAULT NULL," +
		"`dropoff_datetime` datetime DEFAULT NULL," +
		"`passenger_count` int(11) DEFAULT NULL," +
		"`trip_time_in_secs` int(11) DEFAULT NULL," +
		"`trip_distance` double DEFAULT NULL," +
		"`pickup_longitude` double DEFAULT NULL," +
		"`pickup_latitude` double DEFAULT NULL," +
		"`dropoff_longitude` double DEFAULT NULL," +
		"`dropoff_latitude` double DEFAULT NULL" +
		") ENGINE=InnoDB DEFAULT CHARSET=latin1;")

	if err != nil {
		t.Errorf("Failed to create table, %s", err)
		return
	}

	_, err = db.db.Exec("INSERT INTO `cab_trip_data` VALUES ('D7D598CD99978BD012A87A76A7C891B7','82F90D5EFE52FDFD2FDEC3EAD6D5771D','VTS',1,'','2013-12-01 00:13:00','2013-12-01 00:31:00',1,1080,3.79,NULL,NULL,NULL,NULL);")// +
	_, err = db.db.Exec("INSERT INTO `cab_trip_data` VALUES ('5455D5FF2BD94D10B304A15D4B7F2735','177B80B867CEC990DA166BA1D0FCAF82','VTS',1,'','2013-12-01 00:40:00','2013-12-01 00:48:00',6,480,3.2,NULL,NULL,NULL,NULL);")// +
	_, err = db.db.Exec("INSERT INTO `cab_trip_data` VALUES ('9A80FE5419FEA4F44DB8E67F29D84A0F','-73.972794','VTS',1,'-73.995262','2013-12-31 07:39:00','2013-12-31 07:46:00',5,420,2.29,NULL,NULL,NULL,NULL);")// +


	if err != nil {
		t.Errorf("Failed to insert, %s", err)
		return
	}

	trips, err := db.FindTrips(time.Date(2013, 12, 1, 0, 0, 0, 0, time.Local), []string{"D7D598CD99978BD012A87A76A7C891B7"}, true)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(trips))
	assert.Equal(t, "D7D598CD99978BD012A87A76A7C891B7", trips[0].Medallion)
	assert.Equal(t, 1, trips[0].Total)
	assert.Equal(t, "2013-12-1", trips[0].Date)

}
