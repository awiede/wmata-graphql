package wmata_graphql

import (
	"context"
	"github.com/awiede/wmata-go-sdk/wmata"
	"github.com/awiede/wmata-go-sdk/wmata/railpredictions"
	"github.com/pkg/errors"
	"strings"

	"github.com/awiede/wmata-go-sdk/wmata/incidents"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	incidentService        *incidents.Service
	railPredictionsService *railpredictions.Service
}

func NewResolver(apiKey string) *Resolver {
	client := wmata.NewWMATADefaultClient(apiKey)

	incidents := incidents.NewService(client, wmata.JSON)
	railPredictions := railpredictions.NewService(client, wmata.JSON)

	return &Resolver{
		incidentService:        incidents,
		railPredictionsService: railPredictions,
	}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) RailIncident() RailIncidentResolver {
	return &railIncidentResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) BusIncidents(ctx context.Context, route *string) ([]*incidents.BusIncident, error) {
	var busRoute string

	if route != nil {
		busRoute = *route
	}

	response, err := r.incidentService.GetBusIncidents(busRoute)

	if err != nil {
		return nil, err
	}

	busIncidents := make([]*incidents.BusIncident, len(response.BusIncidents))

	for _, busIncident := range response.BusIncidents {
		busIncidents = append(busIncidents, &busIncident)
	}

	return busIncidents, nil

}
func (r *queryResolver) ElevatorIncidents(ctx context.Context, stationID *string) ([]*incidents.ElevatorIncident, error) {
	var railStation string

	if stationID != nil {
		railStation = *stationID
	}

	response, err := r.incidentService.GetOutages(railStation)

	if err != nil {
		return nil, err
	}

	elevatorIncidents := make([]*incidents.ElevatorIncident, len(response.ElevatorIncidents))

	for _, elevatorIncident := range response.ElevatorIncidents {
		elevatorIncidents = append(elevatorIncidents, &elevatorIncident)
	}

	return elevatorIncidents, nil

}
func (r *queryResolver) RailIncidents(ctx context.Context) ([]*incidents.RailIncident, error) {
	response, err := r.incidentService.GetRailIncidents()

	if err != nil {
		return nil, err
	}

	railIncidents := make([]*incidents.RailIncident, len(response.RailIncidents))

	for _, railIncident := range response.RailIncidents {
		railIncidents = append(railIncidents, &railIncident)
	}

	return railIncidents, nil
}

func (r *queryResolver) NextTrains(ctx context.Context, stationCodes []*string) ([]*railpredictions.Train, error) {
	var stations []string

	if stationCodes != nil {
		stations = make([]string, len(stationCodes))

		for _, station := range stationCodes {
			stations = append(stations, *station)
		}
	}

	response, err := r.railPredictionsService.GetNextTrains(stations)

	if err != nil {
		return nil, err
	}

	trains := make([]*railpredictions.Train, len(response.Trains))

	for _, train := range response.Trains {
		trains = append(trains, &train)
	}

	return trains, nil
}

type railIncidentResolver struct{ *Resolver }

func (r *railIncidentResolver) LinesAffected(ctx context.Context, obj *incidents.RailIncident) ([]*string, error) {
	linesSlice := strings.Split(obj.LinesAffected, ",")

	if linesSlice == nil || len(linesSlice) == 0 {
		return nil, errors.New("no lines affected")
	}

	linesAffected := make([]*string, len(linesSlice))

	for _, line := range linesSlice {
		linesAffected = append(linesAffected, &line)
	}

	return linesAffected, nil
}
