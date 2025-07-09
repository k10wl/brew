# Functional Requirements - Kombucha Brewing Tracker

## Problem Statement
Kombucha brewers must track brewing information for multiple jars simultaneously. Current tracking methods fail to provide:
- Reliable data persistence across brewing cycles
- Quick access to jar-specific information during brewing
- Complete historical record of brewing activities per jar
- Data entry speed that doesn't interrupt brewing workflow
- Scalability for managing multiple concurrent batches

System must solve jar identification and data management problems.

## Core Requirements
1. System identifies specific jars uniquely
2. System links brewing data to specific jars
3. System provides rapid access to jar-specific data
4. System stores brewing information persistently
5. System displays complete jar history on demand

---

## FR-1: Jar Identification

### FR-1.1: Unique Jar Identification
System assigns unique identifier to each jar.
- Identifier must be globally unique
- Identifier must be persistent across system updates
- Identifier must never be lost by the system
- System must allow user-friendly names/labels for jars (e.g., "Скубі", "Big Bertha")
- User names must be displayed alongside technical identifiers

### FR-1.2: Physical-Digital Linking
System connects physical jars to digital records.
- System creates downloadable identifiers
- System accepts camera scans and manual entry
- Works on any platform without special hardware

### FR-1.3: Jar Record Persistence
System maintains permanent link between jar identifier and data record.
- Link must survive system updates
- Link must survive data migrations
- Link must provide recovery options for broken associations

---

## FR-2: Data Entry

### FR-2.1: Brewing Recipe Details
System records brewing ingredients and quantities.
- Water amount
- Sugar amount and type
- Tea amount and type
- Additional ingredients

### FR-2.2: Brewing Timeline
System tracks brewing schedule and timing.
- Start date and time
- Estimated harvest time
- Refill schedule
- Event chronology

### FR-2.3: Production Quality
System records post-brewing evaluation.
- Taste rating
- Quality notes
- Brewing success metrics
- Improvement suggestions

### FR-2.4: Data Integrity
System ensures data reliability and consistency.
- System must persist all entered data
- System must validate data before storage
- System must prevent data loss during entry
- System must handle entry failures gracefully
- System must maintain data consistency across sessions

---

## FR-3: Data Access and Management

### FR-3.1: Complete Jar Information Access
System provides comprehensive view of jar brewing data.
- Display all current jar data on scan or selection
- Highlight latest entry and next action

### FR-3.2: Brewing History Tracking
System maintains complete historical record of jar activities.
- Keep and show full jar history
- Include timestamp and session identifier on every record
- Provide filter, search, and sort

### FR-3.3: Record Notes
Brewing records are immutable once saved.
- Records cannot be altered or deleted after submission
- Brewers may append simple notes (e.g., tasting comments, labels) to any record

### FR-3.4: Offline Data Integrity
System protects data when the device goes offline.
- Store entries locally while offline and sync automatically when back online
- Ensure no data loss during offline/online transitions

---

## FR-4: Platform Requirements

### FR-4.1: Browser Support
- Modern browsers
- Mobile & desktop
- No plugins

### FR-4.2: Mobile Optimization
- Touch controls
- One-hand layout
- Large buttons
- Fast entry
- Smart keyboard


### FR-4.3: Offline Operation
- Cache entries
- Offline banner
- Auto sync

### FR-4.4: Performance Requirements
- Fast load
- Quick response
- Instant scan
- Offline ready
- High uptime

---

## FR-5: Data Management

### FR-5.1: Data Storage
System maintains persistent data storage.
- Durability: survives system restarts and failures
- Backup: automated backup mechanisms
- Recovery: data recovery procedures
- Capacity: graceful handling of storage limits
- Retention: minimum 5-year data retention

### FR-5.2: Data Integrity
System ensures data consistency and accuracy.
- Corruption prevention: data validation and checksums
- Relationship integrity: maintains QR code to record links
- Verification: data integrity checking mechanisms

---

## FR-6: Session-Based Jar Management

### FR-6.1: Session Identification
System assigns a unique identifier to each browser session.
- Auto-generated on first visit
- Persisted across refreshes and restarts until cleared
- Requires no personal information
- Displayed on request for troubleshooting and sharing

### FR-6.2: Jar Association
System binds jars and their records to the originating session.
- New jars inherit the current session ID
- Records automatically reference the jar and session
- Jar ownership can be transferred via explicit share/import flow
- Prevent accidental cross-session data mixing

### FR-6.3: Session Sharing
System enables optional sharing of session data.
- Generates share tokens or QR codes that grant access
- Share scopes: read-only or read-write
- Allow revocation of share tokens at any time
- Log all shared edits with timestamp

### FR-6.4: Session Persistence & Security
System safeguards session data and isolation.
- Operates offline first, syncing when online
- Encrypts data at rest and in transit when syncing
- Configurable expiration for inactive sessions
- Isolation guarantees: no data leakage between sessions

---

## FR-7: Notifications and Reminders

### FR-7.1: Brewing Timeline Notifications
System sends reminders for brewing schedule events.
- Harvest time alerts: notify when fermentation period complete
- Refill reminders: alert when continuous brew ready for harvest
- Quality check prompts: remind users to taste and evaluate batches
- Schedule flexibility: user-configurable reminder timing

### FR-7.2: Multi-Jar Coordination
System manages notifications across multiple concurrent batches.
- Jar-specific alerts: separate notifications for each brewing jar
- Batch prioritization: highlight most urgent brewing actions
- Consolidated view: summary of all pending brewing activities
- Smart scheduling: avoid notification overload during busy periods

### FR-7.3: Notification Delivery Methods
System provides multiple notification channels.
- Push notifications: browser and mobile app notifications
- In-app indicators: visual cues within the application

### FR-7.4: Notification Management
System allows user control over notification preferences.
- Customization: user-selectable notification types and timing
- Scheduling: quiet hours and notification frequency settings
- Opt-out: granular control over notification categories
- Snooze functionality: temporary postponement of reminders

### FR-7.5: Context-Aware Alerts
System provides intelligent, situational notifications.
- Weather integration: alerts about temperature effects on brewing
- Seasonal adjustments: brewing time modifications for ambient conditions
- Experience-based: personalized timing based on user's brewing history
- Problem detection: alerts for unusual brewing patterns or potential issues

---

## Implementation Plan (MVP)

### MVP Functional Scope
- FR-1.1: Unique Jar Identification
- FR-1.2: Physical-Digital Linking
- FR-1.4: Jar Record Persistence
- FR-2.1: Brewing Recipe Details
- FR-2.2: Brewing Timeline
- FR-2.3: Production Quality
- FR-3.1: Complete Jar Information Access
- FR-3.2: Brewing History Tracking
- FR-4.1: Browser Support
- FR-4.2: Mobile Optimization
- FR-6.1: Session Identification
- FR-6.2: Jar Association

### MVP Non-Functional Goals
- Offline-first operation for all core flows (create jar, record brew, view history)
- Data persistence across browser restarts and offline periods
- Responsive performance on modern mobile and desktop browsers

### Acceptance Criteria
The MVP is considered complete and ready for release when **all** of the following are true:
1. A user can create a new jar, assign it a friendly name, and obtain a downloadable/scannable identifier (QR code or similar).
2. Scanning or manually entering a jar identifier immediately displays the jar's latest data and full brewing history.
3. Users can record brewing recipe details, timeline events, and production quality notes for any jar without losing data, even if the device goes offline during entry.
4. All data entered offline syncs automatically and correctly when the device reconnects, with no duplication or loss.
5. The application loads in under 2 seconds on a 4G mobile connection and responds to user actions within 300 ms.
6. The interface is fully usable on touch devices with one-handed operation, large tap targets, and adaptive layout.
7. Session isolation is enforced: new jars inherit the current session ID, and jar data is never visible to other sessions unless explicitly shared (future scope).
8. Critical bugs (severity 1) affecting jar creation, data entry, or data retrieval are fixed, and no severity 1 regressions are introduced according to the QA test suite.

## Acceptance Criteria
Implementation complete when all specified requirements met and system enables QR-based brewing data management with superior speed and reliability compared to physical tracking methods. 
