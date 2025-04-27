class API {
    constructor() {
        // Initialize local storage if needed
        if (!localStorage.getItem('queue')) {
            localStorage.setItem('queue', JSON.stringify([]));
        }
        if (!localStorage.getItem('nextId')) {
            localStorage.setItem('nextId', '1');
        }
        if (!localStorage.getItem('currentlyServing')) {
            localStorage.setItem('currentlyServing', '0');
        }
    }

    // Helper methods
    _getQueue() {
        return JSON.parse(localStorage.getItem('queue') || '[]');
    }

    _setQueue(queue) {
        localStorage.setItem('queue', JSON.stringify(queue));
    }

    _getNextId() {
        const id = parseInt(localStorage.getItem('nextId') || '1');
        localStorage.setItem('nextId', (id + 1).toString());
        return id;
    }

    _getCurrentlyServing() {
        return parseInt(localStorage.getItem('currentlyServing') || '0');
    }

    _setCurrentlyServing(id) {
        localStorage.setItem('currentlyServing', id.toString());
    }

    // Simulate network delay
    async _delay() {
        await new Promise(resolve => setTimeout(resolve, 200));
    }

    async joinQueue(data) {
        await this._delay();
        const queue = this._getQueue();
        const id = this._getNextId();
        
        const entry = {
            id,
            firstName: data.firstName,
            lastName: data.lastName,
            email: data.email || '',
            phoneNumber: data.phoneNumber || '',
            joinedAt: new Date().toISOString(),
            status: 'waiting'
        };

        queue.push(entry);
        this._setQueue(queue);

        return { id: entry.id, message: 'Successfully joined the queue' };
    }

    async checkStatus(id) {
        await this._delay();
        const queue = this._getQueue();
        const entry = queue.find(e => e.id === parseInt(id));
        
        if (!entry) {
            throw new Error('Entry not found');
        }

        const position = queue.filter(e => e.status === 'waiting' && e.id < entry.id).length + 1;
        const currentlyServing = this._getCurrentlyServing();

        return {
            status: entry.status,
            position: entry.status === 'waiting' ? position : 0,
            currentlyServing
        };
    }

    async getQueue() {
        await this._delay();
        const queue = this._getQueue();
        const currentlyServing = this._getCurrentlyServing();
        
        return {
            entries: queue,
            currentlyServing
        };
    }

    async notifyNext() {
        await this._delay();
        const queue = this._getQueue();
        const waitingEntries = queue.filter(e => e.status === 'waiting');
        
        if (waitingEntries.length === 0) {
            throw new Error('Queue is empty');
        }

        const nextEntry = waitingEntries[0];
        nextEntry.status = 'serving';
        this._setCurrentlyServing(nextEntry.id);
        this._setQueue(queue);

        return {
            message: 'Next person notified',
            entry: nextEntry
        };
    }

    async markServed(entry) {
        await this._delay();
        const queue = this._getQueue();
        const targetEntry = queue.find(e => e.id === entry.id);
        
        if (!targetEntry) {
            throw new Error('Entry not found');
        }

        targetEntry.status = 'served';
        if (this._getCurrentlyServing() === entry.id) {
            this._setCurrentlyServing(0);
        }
        this._setQueue(queue);

        return {
            message: 'Marked as served',
            entry: targetEntry
        };
    }

    async clearQueue() {
        await this._delay();
        this._setQueue([]);
        this._setCurrentlyServing(0);
        return { message: 'Queue cleared successfully' };
    }
}

// Create a global API instance
const api = new API(); 