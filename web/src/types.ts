export interface User {
  id: string;
  email: string;
  name: string;
  provider: 'google' | 'microsoft';
  created_at: string;
}

export interface Email {
  id: string;
  thread_id: string;
  from: string;
  to: string[];
  subject: string;
  body: string;
  body_html: string;
  date: string;
  is_read: boolean;
  is_starred: boolean;
  labels: string[];
  attachments: string[];
  ai_category?: string;
  is_spam: boolean;
}

export interface EmailSendRequest {
  to: string[];
  cc?: string[];
  bcc?: string[];
  subject: string;
  body: string;
  is_html?: boolean;
}

export interface AIReplyResponse {
  generated_reply: string;
  suggestions: string[];
}

export interface AICategoryResponse {
  category: string;
  confidence: number;
  reasoning: string;
  suggested_labels: string[];
}

export interface AISpamCheckResponse {
  is_spam: boolean;
  confidence: number;
  reasoning: string;
}

export interface CalendarEvent {
  id: string;
  title: string;
  description: string;
  start: string;
  end: string;
  location: string;
  attendees: string[];
  organizer: string;
  is_all_day: boolean;
}

export interface CalendarEventRequest {
  title: string;
  description?: string;
  start: string;
  end: string;
  location?: string;
  attendees?: string[];
  is_all_day?: boolean;
}

export interface AIScheduleResponse {
  suggested_times: string[];
  reasoning: string;
}
