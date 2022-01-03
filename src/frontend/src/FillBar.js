import './FillBar.css';

const green = '#a2ff00';
const yellow = '#ffdd00';
const red = '#ff5e00';

function FillBar(props) {
    let width = props.sum / props.amount * 100
    let color = green;
    if (width > 100) {
        width = 100;
        color = red;
    } else if (Math.abs(props.sum / props.amount) > ((1 / 12) + (props.now / 100))) {
        color = yellow;
    }
    let style = {
        'width': String(width) + "%",
        'backgroundColor': color
    }

    let nowStyle = {
        'left': String(props.now) + "%",
        'width': "2px",
        'height': "2em",
        'position': "relative",
        'backgroundColor': 'black',
        'top': "-2em"
    }

    return (
        <div className="fillBar">
            <div className="filled" style={style}></div>
            <div className="nowLine" style={nowStyle}></div>
        </div>
    );
}

export default FillBar;
