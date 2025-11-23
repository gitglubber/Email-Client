import axios from 'axios';
import {
  User,
  Email,
  EmailSendRequest,
  AIReplyResponse,
  AICategoryResponse,
  AISpamCheckResponse,
  CalendarEvent,
  CalendarEventRequest,
  AIScheduleResponse,
} from './types';

const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080',
  withCredentials: true,
});

// Auth
export const getCurrentUser = async (): Promise<User> => {
  const { data } = await api.get('/auth/me');
  return data;
};

export const logout = async (): Promise<void> => {
  await api.get('/auth/logout');
};

// Emails
export const listEmails = async (
  maxResults: number = 50,
  pageToken?: string
): Promise<{ emails: Email[]; next_token: string }> => {
  const { data } = await api.get('/api/emails', {
    params: { max: maxResults, page_token: pageToken },
  });
  return data;
};

export const getEmail = async (id: string): Promise<Email> => {
  const { data } = await api.get(`/api/emails/${id}`);
  return data;
};

export const sendEmail = async (req: EmailSendRequest): Promise<void> => {
  await api.post('/api/emails/send', req);
};

export const replyToEmail = async (id: string, body: string): Promise<void> => {
  await api.post(`/api/emails/${id}/reply`, { body });
};

export const aiGenerateReply = async (
  id: string,
  context?: string,
  tone?: string
): Promise<AIReplyResponse> => {
  const { data } = await api.post(`/api/emails/${id}/ai-reply`, {
    context,
    tone,
  });
  return data;
};

export const aiCategorizeEmail = async (id: string): Promise<AICategoryResponse> => {
  const { data } = await api.post(`/api/emails/${id}/categorize`);
  return data;
};

export const aiCheckSpam = async (id: string): Promise<AISpamCheckResponse> => {
  const { data } = await api.post(`/api/emails/${id}/spam-check`);
  return data;
};

// Calendar
export const listCalendarEvents = async (
  start?: string,
  end?: string
): Promise<CalendarEvent[]> => {
  const { data } = await api.get('/api/calendar/events', {
    params: { start, end },
  });
  return data.events;
};

export const createCalendarEvent = async (
  req: CalendarEventRequest
): Promise<CalendarEvent> => {
  const { data } = await api.post('/api/calendar/events', req);
  return data;
};

export const aiSuggestScheduling = async (
  description: string,
  attendees?: string[],
  duration?: number
): Promise<AIScheduleResponse> => {
  const { data } = await api.post('/api/calendar/ai-schedule', {
    description,
    attendees,
    duration,
  });
  return data;
};

export default api;
