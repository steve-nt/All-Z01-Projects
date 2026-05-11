// Task 1: Communication functions

let battleCry = (message) => {
  console.log(message.toUpperCase());
};

let secretOrders = (message) => {
  console.log(message.toLowerCase());
};

// Task 2: Robot teamwork functions

let duos = (robot1, robot2) => {
  console.log(robot1 + ' and ' + robot2 + '!');
};

let duosWork = (robot1, robot2, task) => {
  console.log(robot1 + ' and ' + robot2 + ' ' + task + '!');
};

// Task 3: Function that returns a value

let passButter = () => {
  return 'The butter.';
};
