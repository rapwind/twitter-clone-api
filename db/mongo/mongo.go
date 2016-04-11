package mongo

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	mgo "gopkg.in/mgo.v2"
)

const (
	// VERSION ... mgo version
	VERSION = "0.0.1"

	// MONGODB ... type: MongoDB
	MONGODB = "MongoDB"
)

type (
	// Mongo ... Global Mongo
	Mongo struct {
		DataStore interface{}
		Version   string
		Setuped   bool
	}

	// DB ... MongoDB session structure
	DB struct {
		Use       bool
		Dn        string
		Type      string
		DialInfo  *mgo.DialInfo
		Session   *mgo.Session
		Connected bool
	}

	// Collection ... Mongodb#Collection
	Collection struct {
		*mgo.Collection
		Session *mgo.Session `json:"-" bson:"-"`
	}

	// Option ... get options
	Option struct {
		Session bool
		DbName  string
		ColName string
	}
)

// DEBUG ... Debug flag
var DEBUG = false

// Global Mgd instance
var mongo = &Mongo{
	Version:   VERSION,
	DataStore: make(map[string]interface{}),
	Setuped:   false,
}

// Close ... Return back session into pool
func (c *Collection) Close() {
	c.Session.Close()
}

// MongoDB to string
func (d *DB) String() string {

	return fmt.Sprintf("dn=%s, type=%s, connected=%t, addr=%s, database=%s, session=%p",
		d.Dn,
		d.Type,
		d.Connected,
		d.DialInfo.Addrs,
		d.DialInfo.Database,
		d.Session,
	)
}

// Connect ... Connecting to Mongodb
func (d *DB) Connect() error {
	if d.Connected {
		return nil
	}

	session, err := mgo.DialWithInfo(d.DialInfo)

	if err != nil {
		msg := "Failed mongodb connect."
		Debug("%s", msg)
		return errors.New(msg)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	d.Session = session // Original session

	d.Connected = true

	return err
}

// GetSession ... Get a session (Singleton)
func (d *DB) GetSession(makeSession bool) (*mgo.Session, error) {
	if !d.Connected {
		return nil, errors.New("Do not establish a connection with MongoDB. Advance to the Connect() execution")
	}

	if makeSession {
		s, err := d.CopySession()
		return s, err
	}
	// singleton
	return d.Session, nil
}

// CopySession ... Get a new session
func (d *DB) CopySession() (*mgo.Session, error) {
	s, err := d.GetSession(false)
	if err != nil {
		return nil, err
	}
	// Copy(New) session
	return s.Copy(), nil

}

// GetDataBase ... Get Database
func (d *DB) GetDataBase(dbname string, makeSession bool) (*mgo.Database, error) {
	s, err := d.GetSession(makeSession)
	if err != nil {
		return nil, err
	}

	return s.DB(dbname), nil
}

// GetCollection ... Get Collection
func (d *DB) GetCollection(colname string, makeSession bool) (*Collection, error) {
	s, err := d.GetSession(makeSession)
	if err != nil {
		return nil, err
	}

	// Get collection
	c := s.DB(d.DialInfo.Database).C(colname)

	// wrap
	collection := &Collection{
		c,
		s,
	}

	return collection, nil
}

// GetCollectionWithoutErr ... Get Collection without any errors
// if exit == true, shutdown this application immediatry.
func (d *DB) GetCollectionWithoutErr(colname string, makeSession bool, exit bool) *Collection {
	c, err := d.GetCollection(colname, makeSession)
	if err != nil {
		if exit {
			panic(err.Error())
		}
		fmt.Fprintln(os.Stderr, err)
		return nil
	}
	return c
}

// GetDataStore ... Get a datastore
func GetDataStore() (*DB, error) {
	ds := mongo.DataStore

	if ds == nil {
		return nil, errors.New("Datastore not found")
	}

	if ret, ok := ds.(*DB); ok {
		return ret, nil
	}

	return nil, errors.New("Internal data store type is invalid")

}

// Setup ... Setup mongo datastures
func Setup(ds map[string]interface{}, autoconnect bool) error {

	if mongo.Setuped {
		return errors.New("Already setup performed")
	}

	if !ds["Use"].(bool) {
		Debug("Skip the datastore. dn=%s\n", ds["Dn"].(string))
		return nil
	}

	mongodb := &DB{}

	err := mapstructure.Decode(ds, mongodb)

	if err != nil {
		return err
	}

	if autoconnect == true {
		Debug("auto-connecting dn=%s\n", ds["Dn"].(string))
		err := mongodb.Connect()
		if err != nil {
			return err
		}
	}

	mongo.DataStore = mongodb
	Debug("Add the datastore. dn=%s\n", mongodb.Dn)

	mongo.Setuped = true

	return nil
}

// Debug ... Debug output
func Debug(f string, msgs ...string) {
	if DEBUG {
		fmt.Printf(""+f, strings.Join(([]string)(msgs), " "))
	}
}
