class API {
    constructor(baseURL = 'http://localhost:8080') {
        this.baseURL = baseURL;
        this.token = localStorage.getItem('auth_token');
        this.adminKey = localStorage.getItem('admin_key');
    }

    setToken(token) {
        this.token = token;
        localStorage.setItem('auth_token', token);
    }

    setAdminKey(key) {
        this.adminKey = key;
        localStorage.setItem('admin_key', key);
    }

    clearAuth() {
        this.token = null;
        this.adminKey = null;
        localStorage.removeItem('auth_token');
        localStorage.removeItem('admin_key');
    }

    async joinQueue(data) {
        const response = await fetch(`${this.baseURL}/join`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });

        if (!response.ok) {
            throw new Error(`Failed to join queue: ${response.statusText}`);
        }

        const result = await response.json();
        if (result.token) {
            this.setToken(result.token);
        }
        return result;
    }

    async checkStatus(id) {
        if (!this.token) {
            throw new Error('Authentication required');
        }

        const response = await fetch(`${this.baseURL}/status/${id}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${this.token}`,
            },
        });

        if (!response.ok) {
            if (response.status === 401) {
                this.clearAuth();
            }
            throw new Error(`Failed to check status: ${response.statusText}`);
        }

        return response.json();
    }

    async getQueue() {
        if (!this.adminKey) {
            throw new Error('Admin authentication required');
        }

        const response = await fetch(`${this.baseURL}/queue`, {
            method: 'GET',
            headers: {
                'X-API-Key': this.adminKey,
            },
        });

        if (!response.ok) {
            if (response.status === 401) {
                this.clearAuth();
            }
            throw new Error(`Failed to get queue: ${response.statusText}`);
        }

        return response.json();
    }

    async notifyNext() {
        if (!this.adminKey) {
            throw new Error('Admin authentication required');
        }

        const response = await fetch(`${this.baseURL}/next`, {
            method: 'POST',
            headers: {
                'X-API-Key': this.adminKey,
            },
        });

        if (!response.ok) {
            if (response.status === 401) {
                this.clearAuth();
            }
            throw new Error(`Failed to notify next: ${response.statusText}`);
        }

        return response.json();
    }

    async markServed(entry) {
        if (!this.adminKey) {
            throw new Error('Admin authentication required');
        }

        const response = await fetch(`${this.baseURL}/serve`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-API-Key': this.adminKey,
            },
            body: JSON.stringify(entry),
        });

        if (!response.ok) {
            if (response.status === 401) {
                this.clearAuth();
            }
            throw new Error(`Failed to mark as served: ${response.statusText}`);
        }

        return response.json();
    }

    async clearQueue() {
        if (!this.adminKey) {
            throw new Error('Admin authentication required');
        }

        const response = await fetch(`${this.baseURL}/clear`, {
            method: 'POST',
            headers: {
                'X-API-Key': this.adminKey,
            },
        });

        if (!response.ok) {
            if (response.status === 401) {
                this.clearAuth();
            }
            throw new Error(`Failed to clear queue: ${response.statusText}`);
        }

        return response.json();
    }
}

// Create a global API instance
const api = new API(); 