import '@testing-library/jest-dom';
import { fireEvent } from '@testing-library/dom';
import { login, validateToken } from '../js/auth.js';

// Mock fetch for testing API calls
global.fetch = jest.fn();

describe('Authentication', () => {
  beforeEach(() => {
    fetch.mockClear();
    localStorage.clear();
  });

  test('login success returns valid token', async () => {
    const mockToken = 'mock.jwt.token';
    fetch.mockImplementationOnce(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ token: mockToken }),
      })
    );

    const result = await login('testuser@example.com', 'password123');
    expect(result.success).toBe(true);
    expect(result.token).toBe(mockToken);
    expect(localStorage.getItem('token')).toBe(mockToken);
    expect(fetch).toHaveBeenCalledWith('/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: 'testuser@example.com',
        password: 'password123',
      }),
    });
  });

  test('login failure returns error', async () => {
    fetch.mockImplementationOnce(() =>
      Promise.resolve({
        ok: false,
        status: 401,
        json: () => Promise.resolve({ error: 'Invalid credentials' }),
      })
    );

    const result = await login('wrong@example.com', 'wrongpass');
    expect(result.success).toBe(false);
    expect(result.error).toBe('Invalid credentials');
    expect(localStorage.getItem('token')).toBeNull();
  });

  test('validateToken returns true for valid token', async () => {
    const mockToken = 'valid.jwt.token';
    localStorage.setItem('token', mockToken);
    
    fetch.mockImplementationOnce(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ valid: true }),
      })
    );

    const isValid = await validateToken();
    expect(isValid).toBe(true);
    expect(fetch).toHaveBeenCalledWith('/api/validate-token', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${mockToken}`,
      },
    });
  });

  test('validateToken returns false for invalid token', async () => {
    const mockToken = 'invalid.jwt.token';
    localStorage.setItem('token', mockToken);
    
    fetch.mockImplementationOnce(() =>
      Promise.resolve({
        ok: false,
        status: 401,
      })
    );

    const isValid = await validateToken();
    expect(isValid).toBe(false);
  });
}); 