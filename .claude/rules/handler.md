# Handler conventions

Handlers are HTTP-only: parse, validate, delegate, encode. No business logic.

## Constructor

```go
type XHandler struct {
    service   *XService
    validator *validator.Validate
}

func NewXHandler(s *XService, v *validator.Validate) *XHandler {
    return &XHandler{service: s, validator: v}
}
```

## Request flow

Every handler that takes a body follows this exact shape:

```go
func (h *XHandler) Create(w http.ResponseWriter, r *http.Request) {
    var body CreateRequestDTO
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
        return
    }

    if err := h.validator.Struct(body); err != nil {
        shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
        return
    }

    id, err := h.service.Create(r.Context(), body.Field1, body.Field2)
    if err != nil {
        shared_model.WriteErrorResponse(w, err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Header().Set("Location", fmt.Sprintf("/x/%d", id))
    json.NewEncoder(w).Encode(CreateResponseDTO{ID: id, ...})
}
```

Reference: `internal/simp-target/handler.go:35-63`, `internal/user/handler.go:32-59`.

## Path params

```go
idStr := r.PathValue("id")
id, err := strconv.ParseInt(idStr, 10, 64)
if err != nil {
    shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
    return
}
```

Always use `r.PathValue(...)` (Go 1.22+ `ServeMux`), never URL parsing. IDs are `int64`.

## Status codes

- `200 OK` — read (`GET`) returning a body.
- `201 Created` — write returning the created resource. Set `Location: /resource/{id}`.
- `204 No Content` — write that does not return a body (`PUT`, `DELETE`). Do not encode anything.

## Error responses

Never write your own status or message. Always go through `shared_model.WriteErrorResponse(w, err)`. The helper inspects the wrapped sentinel and produces the right status + machine code + message. See [error-handling.md](error-handling.md).

## Content-Type

Always set `w.Header().Set("Content-Type", "application/json")` for handlers that write a JSON body. Do this *before* `WriteHeader`.

## Swagger annotations

Every exported handler MUST carry the full annotation block before the function. Pattern:

```go
// Create creates an X.
// @Summary      Create X
// @Description  Creates an X owned by the current user.
// @Tags         x
// @Accept       json
// @Produce      json
// @Param        x    body      CreateRequestDTO     true  "X data"
// @Success      201  {object}  CreateResponseDTO
// @Failure      400  {object}  model.ErrorResponse  "Invalid request"
// @Failure      401  {object}  model.ErrorResponse  "Unauthorized"
// @Failure      404  {object}  model.ErrorResponse  "Related entity not found"
// @Failure      500  {object}  model.ErrorResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /x [post]
```

- Drop `@Security BearerAuth` on unauthenticated routes (e.g. `/public/...`).
- After adding or changing annotations, regenerate docs: `swag init -g cmd/web/main.go -o cmd/web/docs`.

## DTOs (`dto.go`)

Request DTOs carry `validator` tags; response DTOs do not. Example:

```go
type CreateRequestDTO struct {
    Name        string `json:"name" validate:"required,min=2"`
    Description string `json:"description" validate:"required"`
}

type CreateResponseDTO struct {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}
```

Common tags in this codebase: `required`, `min=N`, `email`. See `internal/user/dto.go:1-8`.
