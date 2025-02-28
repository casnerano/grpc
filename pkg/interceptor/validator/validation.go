package validator

import (
    "errors"

    "google.golang.org/genproto/googleapis/rpc/errdetails"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/protoadapt"
)

const nestedFieldSeparator = "."

type validator interface {
    Validate() error
}

type validatorAll interface {
    ValidateAll() error
}

type validationSingleError interface {
    Field() string
    Reason() string
    Cause() error
    ErrorName() string
}

type validationMultiError interface {
    AllErrors() []error
}

func validation(req interface{}) error {
    switch v := req.(type) {
    case validatorAll:
        return errToGrpcStatus(v.ValidateAll())
    case validator:
        return errToGrpcStatus(v.Validate())
    }
    return nil
}

func extractDetailsFromError(validationErr error, fieldPrefix string) []protoadapt.MessageV1 {
    var details []protoadapt.MessageV1

    switch ve := validationErr.(type) {
    case validationSingleError:
        details = append(details, extractDetailsFromSingleError(ve, fieldPrefix)...)
    case validationMultiError:
        details = append(details, extractDetailsFromMultiError(ve, fieldPrefix)...)
    }

    return details
}

func extractDetailsFromSingleError(singleErr validationSingleError, fieldPrefix string) []protoadapt.MessageV1 {
    if singleErr.Cause() != nil {
        return extractDetailsFromError(
            singleErr.Cause(),
            fieldPrefix+singleErr.Field()+nestedFieldSeparator,
        )
    }

    return []protoadapt.MessageV1{
        &errdetails.BadRequest_FieldViolation{
            Field:       fieldPrefix + singleErr.Field(),
            Description: singleErr.Reason(),
        },
    }
}

func extractDetailsFromMultiError(multiErr validationMultiError, fieldPrefix string) []protoadapt.MessageV1 {
    var details []protoadapt.MessageV1

    for _, err := range multiErr.AllErrors() {
        var singleErr validationSingleError
        if errors.As(err, &singleErr) {
            details = append(details, extractDetailsFromSingleError(singleErr, fieldPrefix)...)
        }
    }

    return details
}

func errToGrpcStatus(err error) error {
    if err == nil {
        return nil
    }

    s := status.New(codes.InvalidArgument, err.Error())

    details := extractDetailsFromError(err, "")
    if len(details) > 0 {
        s, _ = s.WithDetails(details...)
    }

    return s.Err()
}
