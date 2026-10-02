// Mirrors backend/internal/api/dto.go. This file has no behavior of its own
// — it's the shared vocabulary between the typed client and the rest of the app.

export type Operation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'power'
  | 'sqrt'
  | 'percentage';

export interface CalculateRequest {
  operation: Operation;
  operands: number[];
}

export interface CalculateResponse {
  result: number;
}

export interface ArityInfo {
  min: number;
  max: number;
}

export interface OperationInfo {
  name: Operation;
  arity: ArityInfo;
}

export interface OperationsResponse {
  operations: OperationInfo[];
}

export type ErrorCode =
  | 'INVALID_JSON'
  | 'UNKNOWN_FIELD'
  | 'MISSING_FIELD'
  | 'UNKNOWN_OPERATION'
  | 'INVALID_OPERAND_COUNT'
  | 'UNSUPPORTED_MEDIA_TYPE'
  | 'REQUEST_TOO_LARGE'
  | 'DIVISION_BY_ZERO'
  | 'MATH_DOMAIN_ERROR'
  | 'RESULT_OUT_OF_RANGE'
  | 'NOT_FOUND'
  | 'METHOD_NOT_ALLOWED'
  | 'INTERNAL_ERROR';

export interface ErrorBody {
  code: ErrorCode;
  message: string;
}

export interface ErrorResponse {
  error: ErrorBody;
}
