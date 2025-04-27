// Initialize the application when the DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    // Show the join queue section by default
    ui.showSection('joinQueue');

    // Set up auto-refresh for the queue list in admin section
    setInterval(() => {
        if (!document.getElementById('adminSection').classList.contains('hidden')) {
            ui.refreshQueueList();
        }
    }, 5000); // Refresh every 5 seconds if admin section is visible
}); 