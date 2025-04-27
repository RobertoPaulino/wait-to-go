class UI {
    constructor() {
        // Navigation elements
        this.joinQueueBtn = document.getElementById('joinQueueBtn');
        this.checkStatusBtn = document.getElementById('checkStatusBtn');
        this.adminBtn = document.getElementById('adminBtn');

        // Sections
        this.joinQueueSection = document.getElementById('joinQueueSection');
        this.checkStatusSection = document.getElementById('checkStatusSection');
        this.adminSection = document.getElementById('adminSection');

        // Forms
        this.joinForm = document.getElementById('joinForm');
        this.statusForm = document.getElementById('statusForm');

        // Admin elements
        this.nextInQueue = document.getElementById('nextInQueue');
        this.clearQueue = document.getElementById('clearQueue');
        this.queueEntries = document.getElementById('queueEntries');

        // Toast
        this.toast = document.getElementById('toast');

        this.toastTimeout = null;
        
        // Check if admin key is set
        const adminKey = localStorage.getItem('admin_key');
        if (adminKey) {
            api.setAdminKey(adminKey);
            this.showAdminUI();
        } else {
            this.hideAdminUI();
        }

        // Bind event listeners
        this.bindEvents();
    }

    bindEvents() {
        // Navigation
        this.joinQueueBtn.addEventListener('click', () => this.showSection('joinQueue'));
        this.checkStatusBtn.addEventListener('click', () => this.showSection('checkStatus'));
        this.adminBtn.addEventListener('click', () => this.showSection('admin'));

        // Forms
        this.joinForm.addEventListener('submit', this.handleJoinQueue.bind(this));
        this.statusForm.addEventListener('submit', this.handleStatusCheck.bind(this));

        // Admin controls
        this.nextInQueue.addEventListener('click', this.handleNext.bind(this));
        this.clearQueue.addEventListener('click', this.handleClearQueue.bind(this));
    }

    showSection(section) {
        // Hide all sections
        this.joinQueueSection.classList.add('hidden');
        this.checkStatusSection.classList.add('hidden');
        this.adminSection.classList.add('hidden');

        // Remove active class from all buttons
        this.joinQueueBtn.classList.remove('active');
        this.checkStatusBtn.classList.remove('active');
        this.adminBtn.classList.remove('active');

        // Show selected section
        switch (section) {
            case 'joinQueue':
                this.joinQueueSection.classList.remove('hidden');
                this.joinQueueBtn.classList.add('active');
                break;
            case 'checkStatus':
                this.checkStatusSection.classList.remove('hidden');
                this.checkStatusBtn.classList.add('active');
                break;
            case 'admin':
                this.adminSection.classList.remove('hidden');
                this.adminBtn.classList.add('active');
                this.refreshQueueList();
                break;
        }
    }

    showAdminUI() {
        const adminElements = document.querySelectorAll('.admin-only');
        adminElements.forEach(el => el.classList.remove('hidden'));
    }

    hideAdminUI() {
        const adminElements = document.querySelectorAll('.admin-only');
        adminElements.forEach(el => el.classList.add('hidden'));
    }

    showToast(message, isError = false) {
        const toast = document.getElementById('toast');
        toast.textContent = message;
        toast.className = `toast ${isError ? 'error' : 'success'}`;
        
        if (this.toastTimeout) {
            clearTimeout(this.toastTimeout);
        }
        
        this.toastTimeout = setTimeout(() => {
            toast.className = 'toast hidden';
        }, 3000);
    }

    async handleJoinQueue(event) {
        event.preventDefault();
        const formData = new FormData(event.target);
        const data = {
            firstName: formData.get('firstName'),
            lastName: formData.get('lastName'),
            email: formData.get('email'),
            phoneNumber: formData.get('phoneNumber')
        };

        try {
            const result = await api.joinQueue(data);
            this.showToast(`Successfully joined queue. Your ID is: ${result.id}`);
            event.target.reset();
            this.refreshQueueList();
        } catch (error) {
            this.showToast(error.message, true);
        }
    }

    async handleStatusCheck(event) {
        event.preventDefault();
        const id = document.getElementById('queueId').value;

        try {
            const result = await api.checkStatus(id);
            const statusResult = document.getElementById('statusResult');
            const statusMessage = document.getElementById('statusMessage');
            
            statusResult.classList.remove('hidden');
            statusMessage.textContent = `Status: ${result.entry.status}
                Name: ${result.entry.firstName} ${result.entry.lastName}
                Position in Queue: ${result.position}
                Join Time: ${new Date(result.entry.joinTime).toLocaleString()}`;
        } catch (error) {
            if (error.message === 'Authentication required') {
                this.showToast('Please join the queue first to get a token', true);
            } else {
                this.showToast(error.message, true);
            }
        }
    }

    async handleAdminLogin(event) {
        event.preventDefault();
        const apiKey = document.getElementById('adminKey').value;
        
        try {
            api.setAdminKey(apiKey);
            await api.getQueue(); // Test the API key
            this.showToast('Admin login successful');
            this.showAdminUI();
            this.refreshQueueList();
            event.target.reset();
        } catch (error) {
            api.clearAuth();
            this.hideAdminUI();
            this.showToast('Invalid admin key', true);
        }
    }

    async handleAdminLogout() {
        api.clearAuth();
        this.hideAdminUI();
        this.showToast('Admin logged out');
    }

    async handleNext() {
        try {
            await api.notifyNext();
            this.showToast('Next person notified');
            this.refreshQueueList();
        } catch (error) {
            if (error.message === 'Admin authentication required') {
                this.hideAdminUI();
            }
            this.showToast(error.message, true);
        }
    }

    async refreshQueueList() {
        try {
            const entries = await api.getQueue();
            this.queueEntries.innerHTML = '';

            if (!entries || entries.length === 0) {
                this.queueEntries.innerHTML = '<p>No one in queue</p>';
                return;
            }

            entries.forEach(entry => {
                const div = document.createElement('div');
                div.className = 'queue-item';
                div.innerHTML = `
                    <div>
                        <strong>${entry.firstName} ${entry.lastName}</strong>
                        <br>
                        <small>Joined: ${new Date(entry.joinTime).toLocaleString()}</small>
                    </div>
                    <div>
                        <button onclick="ui.handleMarkServed(${entry.id})" class="secondary-btn">
                            Mark Served
                        </button>
                    </div>
                `;
                this.queueEntries.appendChild(div);
            });
        } catch (error) {
            if (error.message === 'Admin authentication required') {
                this.hideAdminUI();
            }
            this.showToast(error.message, true);
        }
    }

    async handleMarkServed(id) {
        try {
            await api.markServed({ id });
            this.showToast('Entry marked as served');
            this.refreshQueueList();
        } catch (error) {
            if (error.message === 'Admin authentication required') {
                this.hideAdminUI();
            }
            this.showToast(error.message, true);
        }
    }

    async handleClearQueue() {
        if (!confirm('Are you sure you want to clear the queue?')) {
            return;
        }

        try {
            await api.clearQueue();
            this.showToast('Queue cleared');
            this.refreshQueueList();
        } catch (error) {
            if (error.message === 'Admin authentication required') {
                this.hideAdminUI();
            }
            this.showToast(error.message, true);
        }
    }
}

// Create a global UI instance
const ui = new UI(); 