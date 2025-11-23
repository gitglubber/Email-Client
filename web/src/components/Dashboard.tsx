import React, { useState, useEffect } from 'react';
import { User, Email, CalendarEvent } from '../types';
import { getCurrentUser, listEmails, listCalendarEvents, logout } from '../api';
import EmailList from './EmailList';
import EmailDetail from './EmailDetail';
import Calendar from './Calendar';
import Sidebar from './Sidebar';

interface DashboardProps {
  onLogout: () => void;
}

const Dashboard: React.FC<DashboardProps> = ({ onLogout }) => {
  const [user, setUser] = useState<User | null>(null);
  const [emails, setEmails] = useState<Email[]>([]);
  const [selectedEmail, setSelectedEmail] = useState<Email | null>(null);
  const [events, setEvents] = useState<CalendarEvent[]>([]);
  const [activeView, setActiveView] = useState<'inbox' | 'calendar'>('inbox');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const userData = await getCurrentUser();
      setUser(userData);

      const emailData = await listEmails(50);
      setEmails(emailData.emails);

      const eventData = await listCalendarEvents();
      setEvents(eventData);
    } catch (error) {
      console.error('Failed to load data:', error);
      onLogout();
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    try {
      await logout();
      onLogout();
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  const refreshEmails = async () => {
    try {
      const emailData = await listEmails(50);
      setEmails(emailData.emails);
    } catch (error) {
      console.error('Failed to refresh emails:', error);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100 flex">
      <Sidebar
        user={user}
        activeView={activeView}
        onViewChange={setActiveView}
        onLogout={handleLogout}
      />

      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header */}
        <header className="bg-white shadow-sm border-b border-gray-200">
          <div className="px-6 py-4 flex items-center justify-between">
            <h1 className="text-2xl font-bold text-gray-900">
              {activeView === 'inbox' ? 'Inbox' : 'Calendar'}
            </h1>
            <div className="flex items-center gap-4">
              <button
                onClick={activeView === 'inbox' ? refreshEmails : loadData}
                className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                title="Refresh"
              >
                <svg
                  className="w-5 h-5 text-gray-600"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                  />
                </svg>
              </button>
            </div>
          </div>
        </header>

        {/* Main Content */}
        <div className="flex-1 overflow-hidden">
          {activeView === 'inbox' ? (
            <div className="h-full flex">
              <EmailList
                emails={emails}
                selectedEmail={selectedEmail}
                onSelectEmail={setSelectedEmail}
              />
              {selectedEmail && (
                <EmailDetail
                  email={selectedEmail}
                  onClose={() => setSelectedEmail(null)}
                  onRefresh={refreshEmails}
                />
              )}
            </div>
          ) : (
            <Calendar events={events} onRefresh={loadData} />
          )}
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
