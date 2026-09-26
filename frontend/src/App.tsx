import Calculator from './components/Calculator'
import SyntaxHelp from './components/SyntaxHelp'

function App() {
  return (
    <main className="flex min-h-screen flex-col items-center px-4">
      <div className="flex w-full flex-1 items-center justify-center py-10">
        <Calculator />
      </div>
      <SyntaxHelp />
    </main>
  )
}

export default App
