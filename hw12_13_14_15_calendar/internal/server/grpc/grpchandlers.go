package grpcserver

import (
	"context"
	"time"

	helpers "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/helpers"
	pb "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (g *GRPCServer) CreateEvent(ctx context.Context, in *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	t := time.Now()
	var message pb.CreateEventResponse
	id, err := g.app.CreateEvent(ctx, in.GetEvent().Title, in.GetEvent().Userid, in.GetEvent().Description, in.GetEvent().Datestart.AsTime(), in.GetEvent().Datestop.AsTime(), in.GetEvent().GetEventmessagetimedelta().AsDuration())
	if err != nil {
		message.Id = 0
		message.Error = err.Error()
	} else {
		message.Id = int32(id)
		message.Error = ""
	}
	logmessage := helpers.StringBuild("[client GRPC: CreateEvent, Request DateTime: ", time.Now().String(), "Time of request work: ", time.Since(t).String())
	g.logg.Info(logmessage)
	return &message, nil
}

func (g *GRPCServer) GetEvent(ctx context.Context, in *pb.GetEventRequest) (*pb.GetEventResponse, error) {
	t := time.Now()
	var message pb.GetEventResponse
	var pbEvent pb.Event
	event, err := g.app.GetEvent(ctx, int(in.GetId()))
	if err != nil {
		message.Error = err.Error()
	} else {
		pbEvent.Id = int32(event.ID)
		pbEvent.Title = event.Title
		pbEvent.Description = event.Description
		pbEvent.Userid = event.UserID
		pbEvent.Datestart = timestamppb.New(event.DateStart)
		pbEvent.Datestop = timestamppb.New(event.DateStop)
		pbEvent.Eventmessagetimedelta = durationpb.New(event.EventMessageTimeDelta)

		message.Event = &pbEvent
		message.Error = ""
	}
	logmessage := helpers.StringBuild("[client GRPC: GetEvent, Request DateTime: ", time.Now().String(), "Time of request work: ", time.Since(t).String())
	g.logg.Info(logmessage)
	return &message, nil
}

func (g *GRPCServer) UpdateEvent(ctx context.Context, in *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	t := time.Now()
	var message pb.UpdateEventResponse
	err := g.app.UpdateEvent(ctx, int(in.GetEvent().Id), in.GetEvent().Title, in.GetEvent().Userid, in.GetEvent().Description, in.GetEvent().Datestart.AsTime(), in.GetEvent().Datestop.AsTime(), in.GetEvent().GetEventmessagetimedelta().AsDuration())
	if err != nil {
		message.Error = err.Error()
	} else {
		message.Error = ""
	}
	logmessage := helpers.StringBuild("[client GRPC: UpdateEvent, Request DateTime: ", time.Now().String(), "Time of request work: ", time.Since(t).String())
	g.logg.Info(logmessage)
	return &message, nil
}

func (g *GRPCServer) DeleteEvent(ctx context.Context, in *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	t := time.Now()
	var message pb.DeleteEventResponse
	err := g.app.DeleteEvent(ctx, int(in.GetId()))
	if err != nil {
		message.Error = err.Error()
	} else {
		message.Error = ""
	}
	logmessage := helpers.StringBuild("[client GRPC: DeleteEvent, Request DateTime: ", time.Now().String(), "Time of request work: ", time.Since(t).String())
	g.logg.Info(logmessage)
	return &message, nil
}

func (g *GRPCServer) GetEventsOnDayByDay(ctx context.Context, in *pb.GetEventsOnDayRequest) (*pb.GetEventsOnDayResponse, error) {
	t := time.Now()
	var message pb.GetEventsOnDayResponse
	events, err := g.app.GetListEventsonDayByDay(ctx, in.GetDate().AsTime())
	if err != nil {
		message.Error = err.Error()
	} else {
		for _, event := range events {
			pbEvent := pb.Event{}
			pbEvent.Id = int32(event.ID)
			pbEvent.Title = event.Title
			pbEvent.Description = event.Description
			pbEvent.Userid = event.UserID
			pbEvent.Datestart = timestamppb.New(event.DateStart)
			pbEvent.Datestop = timestamppb.New(event.DateStop)
			pbEvent.Eventmessagetimedelta = durationpb.New(event.EventMessageTimeDelta)

			message.Events = append(message.Events, &pbEvent)
		}
		message.Error = ""
	}
	logmessage := helpers.StringBuild("[client GetEventsOnDayByDay: GetEvent, Request DateTime: ", time.Now().String(), "Time of request work: ", time.Since(t).String())
	g.logg.Info(logmessage)
	return &message, nil
}

func (g *GRPCServer) GetEventsOnWeekByDay(ctx context.Context, in *pb.GetEventsOnDayRequest) (*pb.GetEventsOnDayResponse, error) {
	t := time.Now()
	var message pb.GetEventsOnDayResponse
	events, err := g.app.GetListEventsOnWeekByDay(ctx, in.GetDate().AsTime())
	if err != nil {
		message.Error = err.Error()
	} else {
		for _, event := range events {
			pbEvent := pb.Event{}
			pbEvent.Id = int32(event.ID)
			pbEvent.Title = event.Title
			pbEvent.Description = event.Description
			pbEvent.Userid = event.UserID
			pbEvent.Datestart = timestamppb.New(event.DateStart)
			pbEvent.Datestop = timestamppb.New(event.DateStop)
			pbEvent.Eventmessagetimedelta = durationpb.New(event.EventMessageTimeDelta)

			message.Events = append(message.Events, &pbEvent)
		}
		message.Error = ""
	}
	logmessage := helpers.StringBuild("[client GetEventsOnWeekByDay: GetEvent, Request DateTime: ", time.Now().String(), "Time of request work: ", time.Since(t).String())
	g.logg.Info(logmessage)
	return &message, nil
}

func (g *GRPCServer) GetEventsOnMonthByDay(ctx context.Context, in *pb.GetEventsOnDayRequest) (*pb.GetEventsOnDayResponse, error) {
	t := time.Now()
	var message pb.GetEventsOnDayResponse
	events, err := g.app.GetListEventsOnMonthByDay(ctx, in.GetDate().AsTime())
	if err != nil {
		message.Error = err.Error()
	} else {
		for _, event := range events {
			pbEvent := pb.Event{}
			pbEvent.Id = int32(event.ID)
			pbEvent.Title = event.Title
			pbEvent.Description = event.Description
			pbEvent.Userid = event.UserID
			pbEvent.Datestart = timestamppb.New(event.DateStart)
			pbEvent.Datestop = timestamppb.New(event.DateStop)
			pbEvent.Eventmessagetimedelta = durationpb.New(event.EventMessageTimeDelta)

			message.Events = append(message.Events, &pbEvent)

		}
		message.Error = ""
	}
	logmessage := helpers.StringBuild("[client GetEventsOnMonthByDay: GetEvent, Request DateTime: ", time.Now().String(), "Time of request work: ", time.Since(t).String())
	g.logg.Info(logmessage)
	return &message, nil
}
