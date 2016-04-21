package service

import (
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/env"
	"fmt"
	"sync"
	"gopkg.in/mgo.v2"
	"github.com/techcampman/twitter-d-server/logger"
)

func sendNotificationForUser(u *entity.User, n *entity.Notification) (err error) {

	ss, err := ReadSessionsByUser(u)

	var wg sync.WaitGroup
	finChan := make(chan bool)
	installationsChan := make(chan *entity.Installation, len(ss))

	wg.Add(len(ss))
	go func() {
		wg.Wait()
		finChan <- true
	}()

	for _, s := range ss {
		go func(t entity.Session) {
			defer wg.Done()

			i, err := ReadInstallationByID(s.InstallationID)
			if err != nil {
				if err != mgo.ErrNotFound {
					logger.Error(err)
				}
				return
			}
			installationsChan <- i
		}(s)
	}
	LOOP:
	for {
		select {
		case <-finChan:
			break LOOP
		case i := <-installationsChan:
			env.GetPushMessage().Send(n.Text, i.ClientType, i.ArnEndpoint)
		}
	}

	fmt.Println(ss)

	return
}
