package email

import (
	"context"
	"fmt"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/storage/database"
	"github.com/sirupsen/logrus"
)

// InboundMail - holds connections to all the inbound mail accounrts
type InboundMail struct {
	client     *client.Client
	loggedOut  bool
	info       *model.InboudMail
	updateChan chan client.Update
	db         database.Database
}

// Init - create connection to the inbount mail
func Init(db database.Database) (*InboundMail, error) {

	ctx := context.Background()
	ctx, cancelFunc := context.WithDeadline(ctx, time.Now().Add(2*time.Second)) //Expires the context when is more than 2 Seconds
	defer cancelFunc()
	inboundMail := InboundMail{}
	inboundMail.db = db
	info, err := db.GetInboundMail(ctx)
	if err != nil {
		return nil, err
	}
	inboundMail.updateChan = make(chan client.Update, 15)
	inboundMail.info = info
	// Attempt to connect to the IMAP Server

	return &inboundMail, nil

}

// Close - logs the user out the email and closes all the channels
func (i *InboundMail) Close() {
	i.client.Logout()
	close(i.updateChan)
}

// Listen - Connects to the IMAP server and listens for new mails
func (i *InboundMail) Listen() {
	logger := logrus.WithField("func", "InboundMail -> Listen()")
	//Connect
	err := i.connect()
	if err != nil {
		// Do not bother polling
		return
	}
	//get message when client gets logged out
	go func() {
		for {
			select {
			case <-i.client.LoggedOut():
				logger.Print("Email just logged out")
				i.loggedOut = true
			}
		}
	}()

	// listen for updates from the client
	/*
	   go func() {

	   		for update := range i.updateChan {
	   			switch update.(type) {
	   			case *client.StatusUpdate:
	   				statusUpdate := update.(*client.StatusUpdate)
	   				log.Printf("just recevied a status update\n%+v \n", statusUpdate)
	   				// log.Println("Just Recived => StatusUpdate Update")
	   				break
	   			case *client.MailboxUpdate:
	   				mailboxUpdate := update.(*client.MailboxUpdate)
	   				log.Printf("just recevied a mailbox update \n%+v \n", mailboxUpdate)
	   				// log.Println("Just Recived => MailboxUpdate Update")
	   				break
	   			case *client.ExpungeUpdate:
	   				expungeUpdate := update.(*client.ExpungeUpdate)
	   				log.Printf("just recevied an expnge update \n%+v \n", expungeUpdate)
	   				// log.Println("Just Recived => ExpungeUpdate Update")
	   				break
	   			case *client.MessageUpdate:
	   				messageUpdate := update.(*client.MessageUpdate)
	   				log.Printf("just recevied a Message update \n%+v \n", messageUpdate)
	   				// log.Println("Just Recived => MessageUpdate Update")
	   				break
	   			}
	   		}
	   	}()
	*/
	// Now fetch mails for the poll period

	// go func(){
	for {
		i.FetchNewMails()
		time.Sleep(time.Duration(*i.info.PollPeriod) * time.Minute)
	}
	// }()

}

// FetchNewMails - checks the imap server for new mails
func (i *InboundMail) FetchNewMails() {
	// Now update the last sequnce on the database
	ctx := context.Background()
	ctx, cancelFunc := context.WithDeadline(ctx, time.Now().Add(3*time.Second)) //Expires the context when is more than 2 Seconds
	defer cancelFunc()

	logger := logrus.WithField("func", "InboundMail -> FetchNewMails()")
	// Select INBOX

	//Check if the client is stil logged in
	if i.loggedOut {
		err := i.connect()
		if err != nil {
			logger.WithField("Err", err.Error()).Info("Reconnecting to the mail server")
		}
	}
	mbox, err := i.client.Select(*i.info.Mailbox, false)
	if err != nil {
		errMessage := fmt.Sprintf("Selecting mailbox mails error: %s", err.Error())
		i.info.Status = func() *string { s := errMessage; return &s }()
		i.db.UpdateInboudMail(ctx, i.info)
		logger.WithField("Err", err.Error()).Info("Selecting the mailbox")

	}

	// check if there is a new mail after the last one
	if mbox.Messages == uint32(*i.info.LastSeq) {
		return
	}

	// LastSeq represents the seq  of the last mail the server is aware of + 1 because the sequence is zero indexed
	from := uint32(*i.info.LastSeq + 1) //uint32(1)
	to := mbox.Messages

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)
	messages := make(chan *imap.Message, 20)
	done := make(chan error, 1)
	go func() {
		done <- i.client.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	for msg := range messages {
		// TODO: Create ticket using certain criteria (Mostyly using the mail subjects)
		/* 		println("\n\n\n\n")
		   		log.Println("* Subject		: " + msg.Envelope.Subject)
		   		log.Println("* InReplyTo	: " + msg.Envelope.InReplyTo)
		   		log.Println("* MessageId	: " + msg.Envelope.MessageId)
		   		log.Printf("* Seq			: %+v\n", msg.SeqNum)
		   		log.Printf("* Flags			: %+v\n", msg.Flags)
		   		log.Printf("* UID			: %+v\n", msg.Uid)
		   		log.Printf("* Size			: %+v\n", msg.Size)
		   		log.Printf("* InternalDate	: %+v", msg.InternalDate)
		*/
		//log.Printf("\n\n* Envelope %v\n\n\n", msg.Envelope)
		// update the new email seqID
		if msg.SeqNum > uint32(*i.info.LastSeq) {
			i.info.LastSeq = func() *int { u := int(msg.SeqNum); return &u }()
		}
	}


	i.db.UpdateInboudMail(ctx, i.info)

	if err := <-done; err != nil {
		errMessage := fmt.Sprintf("Fetcing mails error: %s", err.Error())
		i.info.Status = func() *string { s := errMessage; return &s }()
		i.db.UpdateInboudMail(ctx, i.info)
		logger.WithField("Err", err.Error()).Info("Fetching new mails")
	}

}

func (i *InboundMail) connect() error {

	ctx := context.Background()
	ctx, cancelFunc := context.WithDeadline(ctx, time.Now().Add(3*time.Second)) //Expires the context when is more than 2 Seconds
	defer cancelFunc()

	logger := logrus.WithField("func", "InboundMail -> connect()")
	var c *client.Client
	var err error
	address := fmt.Sprintf("%s:%d", *i.info.Address, *i.info.Port)
	if *i.info.Secured {
		c, err = client.DialTLS(address, nil)
	} else {
		c, err = client.Dial(address)
	}
	if err != nil {
		logger.Error(err)
		errMessage := fmt.Sprintf("Connection error: %s", err.Error())
		i.info.Status = func() *string { s := errMessage; return &s }()
		i.db.UpdateInboudMail(ctx, i.info)
		return err
	}

	logger.Info("Mail Server Connected; Now logging in...")
	// Login
	if err := c.Login(*i.info.EmailUser, *i.info.EmailSecret); err != nil {
		logger.Error(err)
		errMessage := fmt.Sprintf("Connection error: %s", err.Error())
		i.info.Status = func() *string { s := errMessage; return &s }()
		i.db.UpdateInboudMail(ctx, i.info)
		return err
	}
	i.client = c
	//pass the update channel to the client connection
	i.client.Updates = i.updateChan
	i.loggedOut = false

	// Update the mail box status

	i.info.Status = func() *string { s := "Connected"; return &s }()
	i.db.UpdateInboudMail(ctx, i.info)
	return nil
}
