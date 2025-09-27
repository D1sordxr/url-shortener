package handler

import "context"

type Handler struct {
	/* useCase */
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetAnalyticsAlias(
	ctx context.Context,
	request GetAnalyticsAliasRequestObject,
) (GetAnalyticsAliasResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) GetSAlias(
	ctx context.Context,
	request GetSAliasRequestObject,
) (GetSAliasResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) PostShorten(
	ctx context.Context,
	request PostShortenRequestObject,
) (PostShortenResponseObject, error) {
	//TODO implement me
	panic("implement me")
}
