import React, { useState } from 'react';
import { Email, AIReplyResponse } from '../types';
import {
  aiGenerateReply,
  aiCategorizeEmail,
  aiCheckSpam,
  replyToEmail,
} from '../api';
import { format } from 'date-fns';

interface EmailDetailProps {
  email: Email;
  onClose: () => void;
  onRefresh: () => void;
}

const EmailDetail: React.FC<EmailDetailProps> = ({
  email,
  onClose,
  onRefresh,
}) => {
  const [showReply, setShowReply] = useState(false);
  const [replyText, setReplyText] = useState('');
  const [aiReply, setAiReply] = useState<AIReplyResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [showAIOptions, setShowAIOptions] = useState(false);

  const handleAIReply = async () => {
    setLoading(true);
    try {
      const response = await aiGenerateReply(email.id);
      setAiReply(response);
      setReplyText(response.generated_reply);
      setShowReply(true);
    } catch (error) {
      console.error('Failed to generate AI reply:', error);
      alert('Failed to generate AI reply');
    } finally {
      setLoading(false);
    }
  };

  const handleCategorize = async () => {
    setLoading(true);
    try {
      const category = await aiCategorizeEmail(email.id);
      alert(
        `Category: ${category.category}\nConfidence: ${(
          category.confidence * 100
        ).toFixed(1)}%\nReasoning: ${category.reasoning}`
      );
    } catch (error) {
      console.error('Failed to categorize:', error);
      alert('Failed to categorize email');
    } finally {
      setLoading(false);
    }
  };

  const handleSpamCheck = async () => {
    setLoading(true);
    try {
      const spamCheck = await aiCheckSpam(email.id);
      alert(
        `Is Spam: ${spamCheck.is_spam ? 'Yes' : 'No'}\nConfidence: ${(
          spamCheck.confidence * 100
        ).toFixed(1)}%\nReasoning: ${spamCheck.reasoning}`
      );
    } catch (error) {
      console.error('Failed to check spam:', error);
      alert('Failed to check spam');
    } finally {
      setLoading(false);
    }
  };

  const handleSendReply = async () => {
    if (!replyText.trim()) return;

    setLoading(true);
    try {
      await replyToEmail(email.id, replyText);
      alert('Reply sent successfully!');
      setShowReply(false);
      setReplyText('');
      setAiReply(null);
      onRefresh();
    } catch (error) {
      console.error('Failed to send reply:', error);
      alert('Failed to send reply');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex-1 bg-white flex flex-col overflow-hidden">
      {/* Header */}
      <div className="p-6 border-b border-gray-200 flex items-start justify-between">
        <div className="flex-1">
          <h2 className="text-2xl font-semibold text-gray-900 mb-2">
            {email.subject || '(No subject)'}
          </h2>
          <div className="flex items-center gap-4 text-sm text-gray-600">
            <span className="font-medium">{email.from}</span>
            <span>•</span>
            <span>{format(new Date(email.date), 'PPpp')}</span>
          </div>
        </div>
        <button
          onClick={onClose}
          className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
        >
          <svg
            className="w-6 h-6 text-gray-600"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      </div>

      {/* Action Bar */}
      <div className="px-6 py-3 bg-gray-50 border-b border-gray-200 flex items-center gap-2">
        <button
          onClick={() => setShowReply(!showReply)}
          className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          <svg
            className="w-4 h-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6"
            />
          </svg>
          Reply
        </button>

        <button
          onClick={handleAIReply}
          disabled={loading}
          className="flex items-center gap-2 px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors disabled:opacity-50"
        >
          <svg
            className="w-4 h-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
            />
          </svg>
          AI Reply
        </button>

        <div className="relative">
          <button
            onClick={() => setShowAIOptions(!showAIOptions)}
            className="flex items-center gap-2 px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors"
          >
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
              />
            </svg>
            AI Tools
          </button>

          {showAIOptions && (
            <div className="absolute top-full left-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 z-10">
              <button
                onClick={() => {
                  handleCategorize();
                  setShowAIOptions(false);
                }}
                className="w-full text-left px-4 py-2 hover:bg-gray-50 rounded-t-lg"
              >
                Categorize
              </button>
              <button
                onClick={() => {
                  handleSpamCheck();
                  setShowAIOptions(false);
                }}
                className="w-full text-left px-4 py-2 hover:bg-gray-50 rounded-b-lg"
              >
                Check Spam
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Email Body */}
      <div className="flex-1 overflow-y-auto p-6">
        <div className="max-w-4xl">
          {email.body_html ? (
            <div dangerouslySetInnerHTML={{ __html: email.body_html }} />
          ) : (
            <div className="whitespace-pre-wrap text-gray-800">
              {email.body}
            </div>
          )}
        </div>
      </div>

      {/* Reply Box */}
      {showReply && (
        <div className="border-t border-gray-200 p-6 bg-gray-50">
          <div className="max-w-4xl">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              Reply to {email.from}
            </h3>

            {aiReply && aiReply.suggestions.length > 0 && (
              <div className="mb-4">
                <p className="text-sm text-gray-600 mb-2">
                  Quick suggestions:
                </p>
                <div className="flex flex-wrap gap-2">
                  {aiReply.suggestions.map((suggestion, idx) => (
                    <button
                      key={idx}
                      onClick={() => setReplyText(suggestion)}
                      className="text-sm px-3 py-1 bg-blue-100 text-blue-700 rounded-full hover:bg-blue-200 transition-colors"
                    >
                      {suggestion.substring(0, 50)}...
                    </button>
                  ))}
                </div>
              </div>
            )}

            <textarea
              value={replyText}
              onChange={(e) => setReplyText(e.target.value)}
              className="w-full h-40 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
              placeholder="Type your reply..."
            />

            <div className="flex items-center gap-2 mt-4">
              <button
                onClick={handleSendReply}
                disabled={loading || !replyText.trim()}
                className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50"
              >
                Send Reply
              </button>
              <button
                onClick={() => {
                  setShowReply(false);
                  setReplyText('');
                  setAiReply(null);
                }}
                className="px-6 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default EmailDetail;
