const { default: TimeUtils } = require('../data/timeUtils');

function sums(num1, num2) {
  return num1 + num2;
}

describe.each([
  [0, '1 minute'],
  [55, '1 minute'],
  [61, '1 minute'],
  [119, '1 minute'],
  [125, '2 minutes'],
  [901, '15 minutes'],
  [961, '16 minutes'],
])('test parsing seconds into clean string', (input, exp) => {
  test(`${input} -> ${exp}`, () => {
    expect(TimeUtils.secondsToString(input)).toBe(exp);
  });
});
