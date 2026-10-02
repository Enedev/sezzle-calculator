import { afterEach, describe, expect, it, vi } from 'vitest';
import { ApiError, NetworkError, calculate, getOperations } from './client';

function jsonResponse(status: number, body: unknown): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  } as Response;
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('calculate', () => {
  it('resolves with the result on success', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(200, { result: 5 }));
    vi.stubGlobal('fetch', fetchMock);

    const result = await calculate('add', [2, 3]);

    expect(result).toBe(5);
    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining('/calculate'),
      expect.objectContaining({
        method: 'POST',
        headers: expect.objectContaining({ 'Content-Type': 'application/json' }),
        body: JSON.stringify({ operation: 'add', operands: [2, 3] }),
      }),
    );
  });

  it('rejects with an ApiError carrying the backend error code and message', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse(422, { error: { code: 'DIVISION_BY_ZERO', message: 'division by zero' } }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await expect(calculate('divide', [1, 0])).rejects.toMatchObject({
      name: 'ApiError',
      code: 'DIVISION_BY_ZERO',
      message: 'division by zero',
    });
  });

  it('rejects with a NetworkError when fetch itself throws', async () => {
    const fetchMock = vi.fn().mockRejectedValue(new TypeError('Failed to fetch'));
    vi.stubGlobal('fetch', fetchMock);

    await expect(calculate('add', [1, 2])).rejects.toBeInstanceOf(NetworkError);
  });

  it('rejects with an ApiError when the response body is not valid JSON', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => {
        throw new SyntaxError('Unexpected token');
      },
    } as unknown as Response);
    vi.stubGlobal('fetch', fetchMock);

    await expect(calculate('add', [1, 2])).rejects.toBeInstanceOf(ApiError);
  });
});

describe('getOperations', () => {
  it('resolves with the operations list', async () => {
    const operations = [{ name: 'add', arity: { min: 2, max: 2 } }];
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(200, { operations }));
    vi.stubGlobal('fetch', fetchMock);

    const result = await getOperations();

    expect(result).toEqual(operations);
    expect(fetchMock).toHaveBeenCalledWith(expect.stringContaining('/operations'));
  });

  it('rejects with a NetworkError when fetch itself throws', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new TypeError('Failed to fetch')));
    await expect(getOperations()).rejects.toBeInstanceOf(NetworkError);
  });

  it('rejects with an ApiError when the backend returns an error response', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(jsonResponse(500, { error: { code: 'INTERNAL_ERROR', message: 'boom' } })),
    );
    await expect(getOperations()).rejects.toMatchObject({ name: 'ApiError', code: 'INTERNAL_ERROR' });
  });
});
