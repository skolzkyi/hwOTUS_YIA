//go:build integration
// +build integration

package integrationTests

import (
	"testing"
	"time"
	"net/http"
	"context"
	"os/signal"
	"syscall"
	"flag"
	"fmt"
	"os"
	"bytes"
	"io"
	"strconv"

	"github.com/stretchr/testify/require"
	"encoding/json"
	logger "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/logger"
	storage "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/storage/event"
	
	helpers "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/helpers"
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // for driver
)

var configFilePath string
var mySQL_DB *sql.DB
var config Config
var log *logger.LogWrap

type outputJSON struct {
	Text string
	Code int
}

type EventRawData struct {
	EventMessageTimeDelta int64
	Title                 string
	UserID                string
	Description           string
	DateStart             string
	DateStop              string
	ID                    int
}
type EventAnswer struct {
	Events  []storage.Event
	Message outputJSON
}

func init() {
	flag.StringVar(&configFilePath, "config", "./configs/dc/", "Path to config.env")
}


func TestMain(m *testing.M){
	
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config = NewConfig()
	err := config.Init(configFilePath)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("config: ", config)
	log, err = logger.New(config.Logger.Level)
	if err != nil {
		fmt.Println(err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Info("Integration tests down with error")
			os.Exit(1) //nolint:gocritic
		default:
			mySQL_DB,err = InitAndConnectDB(ctx, log, &config)
			if err != nil {
				log.Error("SQL InitAndConnectDB error: " + err.Error())
				cancel()
			}
			err = createTestEventPool(mySQL_DB)
			if err != nil {
				log.Error("SQL DB createTestEventPool error: " + err.Error())
				cancel()
			}
			log.Info("Integration tests up")
  			exitCode := m.Run()
			log.Info("exitCode:"+strconv.Itoa(exitCode))
			//for{} //debug
			err = cleanAndCloseDatabase(ctx, mySQL_DB)
			if err != nil {
    			cancel()
 			}
			 log.Info("Integration tests down succesful")
  			os.Exit(exitCode)//nolint:gocritic
		}
	}
}

func TestCreateEvent(t *testing.T){
	t.Run("CreateEvent_Positive", func(t *testing.T) {
		url := helpers.StringBuild("http://", config.GetServerURL(), "/Event/")
	
		startDateTime := time.Now().Add(800 * time.Hour)
		startDateTimeStr := startDateTime.Format("2006-01-02 15:04:05")

		stopDateTime := time.Now().Add(805 * time.Hour)
		stopDateTimeStr := stopDateTime.Format("2006-01-02 15:04:05")

   	 	jsonStr := []byte(`{"ID":0,"Title":"Control Event from integration test","UserID":"USER0","Description":"","DateStart":"`+startDateTimeStr+`","DateStop":"`+stopDateTimeStr+`","EventMessageTimeDelta": 28800000}`)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
    	require.NoError(t, err)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer resp.Body.Close()

		answer:=outputJSON{}
		err = json.Unmarshal(respBody, &answer)
		require.NoError(t, err)
		

		require.Equal(t, answer.Text, "OK!")

		ctx, cancel := context.WithTimeout(context.Background(), config.GetDBTimeOut())
		defer cancel()

		stmt := "SELECT title FROM eventsTable WHERE id = ?" 
		row := mySQL_DB.QueryRowContext(ctx, stmt, answer.Code)

		var title string

		err = row.Scan(&title)
		require.NoError(t, err)

		require.Equal(t, title, "Control Event from integration test")

		
	})
	t.Run("CreateEvent_Negative_DateBusy", func(t *testing.T) {
		url := helpers.StringBuild("http://", config.GetServerURL(), "/Event/")
	
		startDateTime := time.Now().Add(time.Hour)
		startDateTimeStr := startDateTime.Format("2006-01-02 15:04:05")

		stopDateTime := time.Now().Add(5 * time.Hour)
		stopDateTimeStr := stopDateTime.Format("2006-01-02 15:04:05")

   	 	jsonStr := []byte(`{"ID":0,"Title":"Control Event from integration test","UserID":"USER0","Description":"","DateStart":"`+startDateTimeStr+`","DateStop":"`+stopDateTimeStr+`","EventMessageTimeDelta": 28800000}`)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
		require.NoError(t, err)
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
	
		answer :=outputJSON{}
		err = json.Unmarshal(respBody, &answer)
		require.NoError(t, err)

		require.Equal(t, answer.Text, "this date busy by other event")
		require.Equal(t, answer.Code, 1)

		
	})
}

func TestGetEventByPeriod(t *testing.T) {
	t.Run("GetEventsOnDayByDay", func(t *testing.T) {
		url := helpers.StringBuild("http://", config.GetServerURL(), "/GetEventsOnDayByDay/")

		curDateTime := time.Now()
		curDateTimeStr := curDateTime.Format("2006-01-02 15:04:05")

		jsonStr := []byte(`{"Date":"`+curDateTimeStr+`"}`)
    	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
		require.NoError(t, err)
   	    req.Header.Set("Content-Type", "application/json")

   		client := &http.Client{}
    	resp, err := client.Do(req)
		require.NoError(t, err)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		answer:=EventAnswer{}
		err = json.Unmarshal(respBody, &answer)
		require.NoError(t, err)

		require.Equal(t, answer.Message.Text, "OK!")
		require.Equal(t, answer.Message.Code, 0)

		resTitle := make(map[string]struct{})
		for _, curEvent := range answer.Events {
			resTitle[curEvent.Title] = struct{}{}
		}
		_, ok := resTitle["test0"]
		require.Equal(t, ok, true)
		_, ok = resTitle["test5"]
		require.Equal(t, ok, true)
		require.Equal(t, len(resTitle), 2)

		resp.Body.Close()
	})
	t.Run("GetEventsOnWeekByDay", func(t *testing.T) {
		url := helpers.StringBuild("http://", config.GetServerURL(), "/GetEventsOnWeekByDay/")
    	
		curDateTime := time.Now()
		curDateTimeStr := curDateTime.Format("2006-01-02 15:04:05")

		jsonStr := []byte(`{"Date":"`+curDateTimeStr+`"}`)
		req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
		require.NoError(t, err)
    	req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
    	require.NoError(t, err)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		answer:=EventAnswer{}
		err = json.Unmarshal(respBody, &answer)
		require.NoError(t, err)

		require.Equal(t, answer.Message.Text, "OK!")
		require.Equal(t, answer.Message.Code, 0)

		resTitle := make(map[string]struct{})
		for _, curEvent := range answer.Events {
			resTitle[curEvent.Title] = struct{}{}
		}
		_, ok := resTitle["test0"]
		require.Equal(t, ok, true)
		_, ok = resTitle["test1"]
		require.Equal(t, ok, true)
		_, ok = resTitle["test2"]
		require.Equal(t, ok, true)
		_, ok = resTitle["test5"]
		require.Equal(t, ok, true)
		require.Equal(t, len(resTitle), 4)

		resp.Body.Close()
	})
	t.Run("GetEventsOnMonthByDay", func(t *testing.T) {
		url := helpers.StringBuild("http://", config.GetServerURL(), "/GetEventsOnMonthByDay/")

    	curDateTime := time.Now()
		curDateTimeStr := curDateTime.Format("2006-01-02 15:04:05")

		jsonStr := []byte(`{"Date":"`+curDateTimeStr+`"}`)
		req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
		require.NoError(t, err)
    	req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
    	require.NoError(t, err)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		answer:=EventAnswer{}
		err = json.Unmarshal(respBody, &answer)
		require.NoError(t, err)

		require.Equal(t, answer.Message.Text, "OK!")
		require.Equal(t, answer.Message.Code, 0)

		resTitle := make(map[string]struct{})
		for _, curEvent := range answer.Events {
			resTitle[curEvent.Title] = struct{}{}
		}
		_, ok := resTitle["test0"]
		require.Equal(t, ok, true)
		_, ok = resTitle["test1"]
		require.Equal(t, ok, true)
		_, ok = resTitle["test2"]
		require.Equal(t, ok, true)
		_, ok = resTitle["test3"]
		require.Equal(t, ok, true)
		_, ok = resTitle["test5"]
		require.Equal(t, ok, true)
		require.Equal(t, len(resTitle), 5)

		resp.Body.Close()
	})
}


func TestSendNotificationByDay(t *testing.T) { 
	t.Run("SendNotification", func(t *testing.T) {
		time.Sleep(60*time.Second)
		stmt := `SELECT title FROM eventsTable  WHERE notifCheck="YES"` 

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
	
		retries := 10
		var resTitle map[string]struct{}
		for i := 0; i < retries; i++ {
			resTitle = make(map[string]struct{})
			
			rows,err := mySQL_DB.QueryContext(ctx, stmt)
			
			require.NoError(t, err)
			assign:=""
			for rows.Next() {
				assign =""
				err = rows.Scan(&assign)
				resTitle[assign] = struct{}{}
			}	
			if err = rows.Err(); err != nil {
				fmt.Println("SQL TestSendNotificationByDay rows error: ",err.Error())
				require.NoError(t, err)
			}
			rows.Close()
			if len(resTitle) < 1 {
				time.Sleep(10*time.Second)
				continue
			}
			break
		}
		_, ok := resTitle["test0"]
		require.Equal(t, ok, true)
		_, ok = resTitle["test5"]
		require.Equal(t, ok, true)
		require.Equal(t, len(resTitle), 2)
	})
}



func  InitAndConnectDB(ctx context.Context, logger storage.Logger, config storage.Config) (*sql.DB, error){
	select {
	case <-ctx.Done():
		return nil,storage.ErrStorageTimeout
	default:
		defer recover()
		var err error
		dsn := helpers.StringBuild(config.GetDBUser(), ":", config.GetDBPassword(), "@tcp(",config.GetDBAddress(),":",config.GetDBPort(),")/", config.GetDBName(), "?parseTime=true") //nolint:lll
		
		
		mySQL_DBinn, err := sql.Open("mysql", dsn)
		if err != nil {
			logger.Error("SQL open error: " + err.Error())
			return nil, err
		}
		
		mySQL_DBinn.SetConnMaxLifetime(config.GetDBConnMaxLifetime())
		mySQL_DBinn.SetMaxOpenConns(config.GetDBMaxOpenConns())
		mySQL_DBinn.SetMaxIdleConns(config.GetDBMaxIdleConns())
		
		err = mySQL_DBinn.PingContext(ctx)
		if err != nil {
			logger.Error("SQL DB ping error: " + err.Error())
			return nil, err
		}
		
		return mySQL_DBinn, nil
	}
}

func cleanAndCloseDatabase(ctx context.Context, mySQL_DB *sql.DB) error {
	stmt := "TRUNCATE TABLE OTUSFinalLab.eventsTable"

	_, err := mySQL_DB.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}

	err = mySQL_DB.Close()
	
	return err
}

func createTestEventPool(mySQL_DB *sql.DB) error{
	std := time.Now().Add(time.Hour)
	emtd := 8 * time.Hour
	tempEmtd:=int64(emtd)/1000000
	tx, err := mySQL_DB.Begin()
	if err != nil {
		return err
	}

	stmtPrep, err := tx.Prepare("INSERT INTO eventsTable(title , userID, description , dateStart, dateStop, eventMessageTimeDelta) VALUES (?,?,?,?,?,?)") //nolint:lll,nolintlint
	
	_, err = stmtPrep.Exec("test0", "USER0", "base event",std, std.Add(4*time.Hour), tempEmtd)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmtPrep.Exec("test1", "USER0", "+5days",std.Add(120*time.Hour), std.Add(124*time.Hour), tempEmtd) //nolint:lll,nolintlint
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmtPrep.Exec("test2", "USER0", "+6 days end date after week",std.Add(144*time.Hour), std.Add(150*time.Hour), tempEmtd) //nolint:lll,nolintlint
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmtPrep.Exec("test3", "USER0", "+8 days - next week",std.Add(192*time.Hour), std.Add(200*time.Hour), tempEmtd) //nolint:lll,nolintlint
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmtPrep.Exec("test4", "USER0", "start in before week and end in cur week",std.Add(-48*time.Hour), std.Add(-6*time.Hour), tempEmtd) //nolint:lll,nolintlint
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmtPrep.Exec("test5", "USER0", "in this day",std.Add(-3*time.Hour), std.Add(-2*time.Hour), tempEmtd) //nolint:lll,nolintlint
	if err != nil {
		tx.Rollback()
		return err
	}
	
	err = tx.Commit()
		
	return err
}