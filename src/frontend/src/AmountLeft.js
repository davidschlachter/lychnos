import { useTheme } from '@mui/material/styles';

const red = '#ff5e00';
const green = '#a2ff00';

function AmountLeft(props) {
    const theme = useTheme();

    let timeSpent;
    if (props.timeSpent > 100) {
        timeSpent = 100;
    } else {
        timeSpent = props.timeSpent
    }

    let left;
    if (timeSpent > 91.667) {
        // From the last month onwards, don't keep increasing the amount with time.
        left = Math.round(props.amount - props.sum)
    } else {
        left = Math.round((props.amount - props.sum) / (12 * (1 - (timeSpent / 100))))
        if (left === 0) {
            left = 1;
        }
    }
    let color = theme.palette.text.primary
    let sign = 1
    if (Math.sign(left) !== Math.sign(props.amount)) {
        if (props.amount > 0) {
            color = green // Celebrate income category greater than the target
        } else if (Math.abs(left) > 10) {
            color = red
        }
        sign = -1
    }

    // If the amount is very close to zero, round it to zero.
    if (left < 2 && left > -2) {
        left = 0;
    }

    // Don't display 'amounts left' for exceeded income categories.
    const displayString = color == green ? "✅" : sign * Math.abs(left)

    return (
        <span style={{
            "color": color
        }}>{displayString}</span>
    );
}

export default AmountLeft;
