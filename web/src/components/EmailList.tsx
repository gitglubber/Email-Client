import React from 'react';
import { Email } from '../types';
import { formatDistanceToNow } from 'date-fns';

interface EmailListProps {
  emails: Email[];
  selectedEmail: Email | null;
  onSelectEmail: (email: Email) => void;
}

const EmailList: React.FC<EmailListProps> = ({
  emails,
  selectedEmail,
  onSelectEmail,
}) => {
  const getCategoryColor = (category?: string) => {
    const colors: { [key: string]: string } = {
      Work: 'bg-blue-100 text-blue-800',
      Personal: 'bg-green-100 text-green-800',
      Finance: 'bg-yellow-100 text-yellow-800',
      Shopping: 'bg-purple-100 text-purple-800',
      Social: 'bg-pink-100 text-pink-800',
      Important: 'bg-red-100 text-red-800',
    };
    return colors[category || ''] || 'bg-gray-100 text-gray-800';
  };

  return (
    <div className="w-96 border-r border-gray-200 bg-white overflow-y-auto">
      <div className="divide-y divide-gray-100">
        {emails.length === 0 ? (
          <div className="p-8 text-center text-gray-500">
            <svg
              className="w-16 h-16 mx-auto mb-4 text-gray-300"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
              />
            </svg>
            <p>No emails found</p>
          </div>
        ) : (
          emails.map((email) => (
            <div
              key={email.id}
              onClick={() => onSelectEmail(email)}
              className={`p-4 cursor-pointer transition-colors ${
                selectedEmail?.id === email.id
                  ? 'bg-blue-50 border-l-4 border-blue-500'
                  : 'hover:bg-gray-50 border-l-4 border-transparent'
              } ${!email.is_read ? 'font-semibold' : ''}`}
            >
              <div className="flex items-start justify-between mb-1">
                <span
                  className={`text-sm truncate flex-1 ${
                    email.is_read ? 'text-gray-600' : 'text-gray-900'
                  }`}
                >
                  {email.from}
                </span>
                <span className="text-xs text-gray-500 ml-2 whitespace-nowrap">
                  {formatDistanceToNow(new Date(email.date), {
                    addSuffix: true,
                  })}
                </span>
              </div>

              <div className="flex items-center gap-2 mb-2">
                <h3
                  className={`text-sm truncate flex-1 ${
                    email.is_read ? 'text-gray-700' : 'text-gray-900'
                  }`}
                >
                  {email.subject || '(No subject)'}
                </h3>
                {email.is_starred && (
                  <svg
                    className="w-4 h-4 text-yellow-500 flex-shrink-0"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                  </svg>
                )}
              </div>

              <p className="text-xs text-gray-500 truncate mb-2">
                {email.body.substring(0, 100)}
              </p>

              {email.ai_category && (
                <span
                  className={`inline-block px-2 py-1 text-xs rounded-full ${getCategoryColor(
                    email.ai_category
                  )}`}
                >
                  {email.ai_category}
                </span>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default EmailList;
