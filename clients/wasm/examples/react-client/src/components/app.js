import React from 'react'

export default class App extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      isLoading: true
    }
  }
  componentDidMount() {
	WebAssembly.instantiateStreaming(fetch("http://localhost:8080/client/main.wasm"), go.importObject).then(async (result) => {
    go.run(result.instance)
    this.setState({ isLoading: false })

    const obj = connect()
    obj.onMessage = (mes) => {console.log("js callback" + mes)}
    console.log(obj.do())
	});
  }
  render() {
    return (
      this.state.isLoading ? 
      <div>Loading</div> :  
      <div>Go WASM React Example </div>
    )
  }
}