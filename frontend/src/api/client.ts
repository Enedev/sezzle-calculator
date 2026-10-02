// The only module in this app that knows the HTTP contract (endpoints,
// request/response shapes, status codes). Everything else talks to the
// backend only through calculate()/getOperations() below.
import type {
  CalculateResponse,
  ErrorCode,
  ErrorResponse,
  Operation,
  OperationInfo,
  OperationsResponse,
} from '../types/api';

const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api/v1';

export class ApiError extends Error {
  code: ErrorCode;

  constructor(code: ErrorCode, message: string) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
  }
}

export class NetworkError extends Error {
  constructor(message = 'Network request failed') {
    super(message);
    this.name = 'NetworkError';
  }
}

async function parseJSON<T>(response: Response): Promise<T> {
  try {
    return (await response.json()) as T;
  } catch {
    throw new ApiError('INTERNAL_ERROR', 'The server returned an unexpected response.');
  }
}

export async function calculate(operation: Operation, operands: number[]): Promise<number> {
  let response: Response;
  try {
    response = await fetch(`${BASE_URL}/calculate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ operation, operands }),
    });
  } catch {
    throw new NetworkError();
  }

  if (!response.ok) {
    const body = await parseJSON<ErrorResponse>(response);
    throw new ApiError(body.error.code, body.error.message);
  }

  const body = await parseJSON<CalculateResponse>(response);
  return body.result;
}

export async function getOperations(): Promise<OperationInfo[]> {
  let response: Response;
  try {
    response = await fetch(`${BASE_URL}/operations`);
  } catch {
    throw new NetworkError();
  }

  if (!response.ok) {
    const body = await parseJSON<ErrorResponse>(response);
    throw new ApiError(body.error.code, body.error.message);
  }

  const body = await parseJSON<OperationsResponse>(response);
  return body.operations;
}
