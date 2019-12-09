"use strict";

class Keyboard extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            input: '',
            translated: [], 
            translatedText: [],
            options: [],
            highlightedOption: -1}
            this.getTranslation = this.getTranslation.bind(this),
            this.handleInput = this.handleInput.bind(this);
            this.handleEnter = this.handleEnter.bind(this);
            this.handleNavigation = this.handleNavigation.bind(this);
            this.handleBackspace = this.handleBackspace.bind(this);
            this.selectOption = this.selectOption.bind(this);
    }

    handleInput(input){
        const str = arrayToString(this.state.translated)
        const length = str.length;
        const newInput = input.substring(length);
        this.getTranslation(newInput);
    }

    getTranslation(input) {
        fetch("/translations?input=" + input)
            .then(res => res.json())
            .then(
                (result => {
                    this.setState({
                        input: input,
                        options: ((result.chars == null) ? [] : result.chars), 
                        highlightedOption: ((result.chars == null) ? -1 : 0)
                    });
                })
            )
    }

    selectOption(value){

        this.setState(state => ({
            input: '',
            translated: state.translated.concat(value),
            translatedText: state.translatedText.concat(state.input),
            options: [], 
            highlightedOption: -1
        }));
    }
    
    handleEnter() {
        if (this.state.highlightedOption >= 0) {
            this.selectOption(this.state.options[this.state.highlightedOption]);
        }
    }

    handleBackspace() {

        //if input length is 1 (before backspace) then remove translations (i.e. don't bother translating empty string)
        if (this.state.input.length === 1) {
            this.setState({
                input: '',
                options: [], 
                highlightedOption: -1
            })
        } else if (this.state.input.length === 0) {

            //if we have translations then reverse translation of last translated character and call translate
            //otherwise do nothing
            if (this.state.translated.length > 0) {
                this.state.translated.pop();

                this.setState(state => ({
                    //input: (state.translatedText.length > 0) ? state.translatedText.pop() : '',
                    input: state.translatedText.pop(),
                    translated: state.translated,
                    translatedText: state.translatedText,
                    options: [], 
                    highlightedOption: -1
                }), () => this.getTranslation(this.state.input));}
            //else just remove last input
        } else {
            this.setState(state => ({
                input: state.input.substring(0, state.input.length-1)
            }),() => this.getTranslation(this.state.input));
        }
    }

    handleNavigation(direction) {
        
        const lenOptions = this.state.options.length;
        var highlightedOption = this.state.highlightedOption;

        if (direction === "left") {
            highlightedOption = (highlightedOption - 1 < 0) ? lenOptions - 1 : highlightedOption - 1
        } else if (direction === "right") {
            highlightedOption = (highlightedOption +1 === lenOptions) ? 0 : highlightedOption + 1
        }

        //alert(highlightedOption);

        if (this.state.highlightedOption >= 0) {
            this.setState({
                highlightedOption: highlightedOption
            });
        }
    }

    render() {

        const highlight = (this.state.highlightedOption > 0) ? this.state.highlightedOption : 0;

        //using index is last resort. use a key like ascii or frequency
        const optionsList = this.state.options.map((option, index) =>
            <Option key={index} className={(index===highlight) ? "highlight" : "plain"} 
                char={option} optionsClick={this.selectOption}/>)
        
        const text = arrayToString(this.state.translated) + this.state.input;

        return (
            <div id="page">
                <div className="header">
                    <div className="headertext">
                        <h1>鍵盤中文拼音</h1>
                        <p>Jianpan is a pinyin keyboard</p>
                    </div>
                    
                </div>
                <div className="row">
                    <div className="column side">
                        {/* <h2>something at the side</h2> */}
                    </div>
                    <div id="keyboard" className="column middle">
                        <h3>Traditional Chinese</h3>
                        <div id="options">{optionsList}</div>
                        <Input keyPress={this.handleInput} keyDown={this.handleEnter} move={this.handleNavigation} 
                            backspace={this.handleBackspace} value={text} />
                    </div>
                    <div className="column side">
                        {/* <h2>something at the side</h2> */}
                    </div>
                </div>
                {/* <div className="footer"><p>footer text</p></div> */}
            </div>
        )
    }
}

class Input extends React.Component {
    constructor(props) {
        super(props);
        this.handleKeyPress = this.handleKeyPress.bind(this);
        this.handleKeyDown = this.handleKeyDown.bind(this);
    }

    handleKeyPress(event) {
        this.props.keyPress(event.target.value);
    }

    handleKeyDown(event) {
        
        if (event.key === 'Enter') {
            this.props.keyDown();
        } else if (event.key === 'Tab') {
            event.preventDefault();
            (event.shiftKey) ? this.props.move("left") : this.props.move("right");
        } else if (event.key === "ArrowLeft") {
            event.preventDefault();
            this.props.move("left");
        } else if (event.key === "ArrowRight") {
            event.preventDefault();
            this.props.move("right");
        } else if (event.key === 'Backspace') {
            event.preventDefault();
            this.props.backspace();
        }
    }

    componentDidMount() {
        this.refs.inputField.focus();
    }

    render() {
        return (
            <div>
                <input ref="inputField" type="text" /*placeholder="type here" */
                value={this.props.value} onChange={this.handleKeyPress} onKeyDown={this.handleKeyDown} onKeyUp={this.handleKeyUp} />
            </div>
        )
    }
}

class Option extends React.Component {
    constructor(props) {
        super(props);
        this.handleClick = this.handleClick.bind(this);
    }

    handleClick(){
        this.props.optionsClick(this.props.char);
    }

    render() {

        return (
            <button onClick={this.handleClick}
                className={this.props.className}>
                {this.props.char}
            </button>
        );
    }
}

function arrayToString(array) {
    var i;
    var arrayString = '';
    for (i=0; i < array.length; i++) {
        arrayString += array[i];
    }
    return arrayString;
}

ReactDOM.render(
    <Keyboard />,
    document.getElementById('root')
  );