package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/concertLabs/oaf-server/graph/generated"
	"github.com/concertLabs/oaf-server/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, user model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, password *string, email *string, showname *string) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateOrganization(ctx context.Context, organization model.NewOrganization) (*model.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateOrganization(ctx context.Context, id string, name *string, picture *string) (*model.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteOrganization(ctx context.Context, id string) (*model.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateSection(ctx context.Context, section model.NewSection) (*model.Section, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateSection(ctx context.Context, id string, name string) (*model.Section, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteSection(ctx context.Context, id string) (*model.Section, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateSectionMember(ctx context.Context, section string, user string, right *int) (*model.Member, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateSectionMember(ctx context.Context, section string, user string, right int) (*model.Member, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteSectionMember(ctx context.Context, section string, user string) (*model.Member, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateEvent(ctx context.Context, event model.NewEvent) (*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateEvent(ctx context.Context, id string, name *string, description *string, adress *string, start *string, end *string) (*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteEvent(ctx context.Context, id string) (*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateEventAttendee(ctx context.Context, event string, user string, commitment int, comment *string) (*model.Attendee, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateEventAttendee(ctx context.Context, event string, user string, commitment int, comment *string) (*model.Attendee, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteEventAttendee(ctx context.Context, event string, user string) (*model.Attendee, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateEventComment(ctx context.Context, event string, text string) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateEventComment(ctx context.Context, id string, text string) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteEventComment(ctx context.Context, id string) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateInvite(ctx context.Context, invite model.NewInvite) (*model.Invite, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteInvite(ctx context.Context, id string) (*model.Invite, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Organization(ctx context.Context, id string) (*model.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Section(ctx context.Context, id string) (*model.Section, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Member(ctx context.Context, id string) (*model.Member, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Members(ctx context.Context, section *string, user *string, right *int) ([]*model.Member, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Event(ctx context.Context, id string) (*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Events(ctx context.Context, organization *string, start *string, end *string) ([]*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Comment(ctx context.Context, id string) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Comments(ctx context.Context, event string) ([]*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Attendee(ctx context.Context, id string) (*model.Attendee, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Attendees(ctx context.Context, event *string, user *string, commitment *model.Commitment) ([]*model.Attendee, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Invite(ctx context.Context, id string) (*model.Invite, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Invites(ctx context.Context, section *string, user *string) ([]*model.Invite, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
