# Use Cases and Requirements

## Use Cases

### UC1: Multi-factor authentication (MFA) and key rotation

A single user may want to have an account that requires at least two out of
three key to sign for instance.

That use may also want to have an account that isn't tied to a single permanent
key but rather where they can add and remove keys as their needs evolve. The
member keys could be tied to a single device or hardware wallet, or potentially
tied to a custodial service that users other mechanisms to support MFA 

### UC2: Organizational Accounts

A group of people that know each other.

### UC3: DAO's

A group of people that may or may not know each other.

## Requirements

### R1: Groups aren't tied to any permanent key

### R2: Groups are themselves an account

### R3: Keys can be added and removed from a group

### R4: Groups can submit and execute a transaction within a single transaction

### R5: Different keys can be assigned different weights

### R6: Different voting strategies and thresholds can be chosen

### R7: The same group can use different voting stategies for different accounts

### R8: Garbage proposals eventually get cleaned up