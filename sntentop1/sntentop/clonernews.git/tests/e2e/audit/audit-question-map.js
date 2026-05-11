/**
 * @typedef {Object} AuditQuestionMapEntryBase
 * @property {string} id
 * @property {'Functional' | 'General' | 'Bonus'} section
 * @property {string} question
 * @property {readonly string[]} implementingTickets
 */

/**
 * @typedef {AuditQuestionMapEntryBase & ({ testFile: string } | { testFiles: readonly string[] })} AuditQuestionMapEntry
 */

/** @type {readonly AuditQuestionMapEntry[]} */
export const AUDIT_QUESTION_MAP = Object.freeze([
  Object.freeze({
    id: 'AUDIT-F-01',
    section: 'Functional',
    question: 'Does this post open without any errors?',
    implementingTickets: Object.freeze(['TC-3']),
    testFile: 'stories.spec.js',
  }),
  Object.freeze({
    id: 'AUDIT-F-02',
    section: 'Functional',
    question: 'Does this post open without any errors?',
    implementingTickets: Object.freeze(['TC-3']),
    testFile: 'jobs.spec.js',
  }),
  Object.freeze({
    id: 'AUDIT-F-03',
    section: 'Functional',
    question: 'Does this post open without any errors?',
    implementingTickets: Object.freeze(['TC-6']),
    testFile: 'polls.spec.js',
  }),
  Object.freeze({
    id: 'AUDIT-F-04',
    section: 'Functional',
    question: 'Did the posts load without error and without spamming the user?',
    implementingTickets: Object.freeze(['TB-3']),
    testFile: 'load-more.spec.js',
  }),
  Object.freeze({
    id: 'AUDIT-F-05',
    section: 'Functional',
    question: 'Are the comments being displayed in the correct order (from newest to oldest)?',
    implementingTickets: Object.freeze(['TC-4']),
    testFile: 'comments.spec.js',
  }),
  Object.freeze({
    id: 'AUDIT-G-01',
    section: 'General',
    question: 'Does the UI have at least stories, jobs and polls?',
    implementingTickets: Object.freeze(['TB-2', 'TC-3', 'TC-6']),
    testFiles: Object.freeze(['stories.spec.js', 'polls.spec.js']),
  }),
  Object.freeze({
    id: 'AUDIT-G-02',
    section: 'General',
    question: 'Are the posts displayed in the correct order (from newest to oldest)?',
    implementingTickets: Object.freeze(['TB-2']),
    testFile: 'stories.spec.js',
  }),
  Object.freeze({
    id: 'AUDIT-G-03',
    section: 'General',
    question: 'Does each comment present the right parent post?',
    implementingTickets: Object.freeze(['TC-4']),
    testFile: 'comments.spec.js',
  }),
  Object.freeze({
    id: 'AUDIT-G-04',
    section: 'General',
    question: 'Does the UI notify the user when a certain post is updated?',
    implementingTickets: Object.freeze(['TC-5']),
    testFile: 'live-data.spec.js',
  }),
  Object.freeze({
    id: 'AUDIT-G-05',
    section: 'General',
    question:
      'Is the project using throttling to regulate the number of requests (every 5 seconds)?',
    implementingTickets: Object.freeze(['TA-4', 'TC-5']),
    testFile: 'live-data.spec.js',
  }),
  Object.freeze({
    id: 'AUDIT-B-01',
    section: 'Bonus',
    question: '+Does the UI have more types of posts than stories, jobs and polls?',
    implementingTickets: Object.freeze(['TB-2']),
    testFile: 'stories.spec.js',
  }),
  Object.freeze({
    id: 'AUDIT-B-02',
    section: 'Bonus',
    question: '+Have sub-comments (nested comments) been implemented?',
    implementingTickets: Object.freeze(['TC-4']),
    testFile: 'comments.spec.js',
  }),
]);

const duplicateAuditIds = [
  ...new Set(
    AUDIT_QUESTION_MAP.map((entry) => entry.id).filter(
      (id, index, allIds) => allIds.indexOf(id) !== index,
    ),
  ),
].toSorted();

if (duplicateAuditIds.length > 0) {
  throw new Error(
    `Duplicate audit IDs detected in AUDIT_QUESTION_MAP: ${duplicateAuditIds.join(', ')}`,
  );
}

export default AUDIT_QUESTION_MAP;
