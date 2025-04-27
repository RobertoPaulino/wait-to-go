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

        // Bind event listeners
        this.bindEvents();
    }

    bindEvents() {
        // Navigation
        this.joinQueueBtn.addEventListener('click', () => this.showSection('joinQueue'));
        this.checkStatusBtn.addEventListener('click', () => this.showSection('checkStatus'));
        this.adminBtn.addEventListener('click', () => this.showSection('admin'));

        // Forms
        this.joinForm.addEventListener('submit', this.handleJoinSubmit.bind(this));
        this.statusForm.addEventListener('submit', this.handleStatusCheck.bind(this));

        // Admin controls
        this.nextInQueue.addEventListener('click', this.handleNextInQueue.bind(this));
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

    async handleJoinSubmit(event) {
        event.preventDefault();
        const formData = new FormData(this.joinForm);
        const data = {
            firstName: formData.get('firstName'),
            lastName: formData.get('lastName'),
            email: formData.get('email'),
            phoneNumber: formData.get('phoneNumber')
        };

        try {
            const response = await api.joinQueue(data);
            this.showToast(`Successfully joined queue! Your ID is: ${response.id}`);
            this.joinForm.reset();
        } catch (error) {
            this.showToast(error.message, true);
        }
    }

    async handleStatusCheck(event) {
        event.preventDefault();
        const id = document.getElementById('queueId').value;

        try {
            const response = await api.checkStatus(id);
            const statusResult = document.getElementById('statusResult');
            const statusMessage = document.getElementById('statusMessage');
            
            let message = `Status: ${response.status}\n`;
            if (response.status === 'waiting') {
                message += `Position in queue: ${response.position}\n`;
            }
            if (response.currentlyServing) {
                message += `Currently serving: #${response.currentlyServing}`;
            }
            
            statusResult.classList.remove('hidden');
            statusMessage.textContent = message;
        } catch (error) {
            this.showToast(error.message, true);
        }
    }

    async handleNextInQueue() {
        try {
            const response = await api.notifyNext();
            this.showToast(response.message);
            this.refreshQueueList();
        } catch (error) {
            this.showToast(error.message, true);
        }
    }

    async handleClearQueue() {
        if (!confirm('Are you sure you want to clear the queue?')) return;

        try {
            const response = await api.clearQueue();
            this.showToast(response.message);
            this.refreshQueueList();
        } catch (error) {
            this.showToast(error.message, true);
        }
    }

    async refreshQueueList() {
        try {
            const response = await api.getQueue();
            this.queueEntries.innerHTML = '';

            if (!response.entries || response.entries.length === 0) {
                this.queueEntries.innerHTML = '<p>No one in queue</p>';
                return;
            }

            response.entries.forEach(entry => {
                const div = document.createElement('div');
                div.className = 'queue-item';
                div.innerHTML = `
                    <div>
                        <strong>${entry.firstName} ${entry.lastName}</strong>
                        <br>
                        <small>ID: #${entry.id} | Status: ${entry.status}</small>
                        <br>
                        <small>Joined: ${new Date(entry.joinedAt).toLocaleString()}</small>
                    </div>
                    ${entry.status === 'serving' ? `
                    <div>
                        <button onclick="ui.handleMarkServed(${entry.id})" class="secondary-btn">
                            Mark Served
                        </button>
                    </div>
                    ` : ''}
                `;
                this.queueEntries.appendChild(div);
            });
        } catch (error) {
            this.showToast(error.message, true);
        }
    }

    async handleMarkServed(id) {
        try {
            const response = await api.markServed({ id });
            this.showToast(response.message);
            this.refreshQueueList();
        } catch (error) {
            this.showToast(error.message, true);
        }
    }

    showToast(message, isError = false) {
        this.toast.textContent = message;
        this.toast.style.background = isError ? 'var(--danger-color)' : 'var(--text)';
        this.toast.classList.remove('hidden');

        setTimeout(() => {
            this.toast.classList.add('hidden');
        }, 3000);
    }
}

// Create a global UI instance
const ui = new UI(); 