// Task 1: Check if truth is truthy or falsy

// A variable truth has been provided
if (truth) {
  console.log('The truth was spoken.');
} else {
  console.log('Lies !!!!');
}

// Task 2: Check multiple conditions (AND operator)

// A user object with age and activeMembership properties has been provided
// A variable ticket already exists
ticket = 'You cannot benefit from our special promotion.';

if (user.age >= 18 && user.age <= 25 && user.activeMembership) {
  ticket = 'You can benefit from our special promotion.';
}

// Task 3: Check if customer can afford ticket (OR operator)

// A customer object with cash and hasVoucher properties has been provided
// A variable ticketSold already exists
if (customer.cash >= 9.99 || customer.hasVoucher) {
  ticketSold += 1;
}
