import Link from "next/link"
import { Button } from "@/components/ui/Button"
import { BookOpen } from "lucide-react"

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col overflow-hidden bg-background relative">
      {/* Background Gradients */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-[600px] bg-gradient-to-b from-indigo-50/80 to-transparent dark:from-indigo-950/30 z-0 pointer-events-none" />
      <div className="absolute top-[10%] left-[10%] h-[400px] w-[400px] rounded-full bg-violet-500/10 blur-[100px] z-0" />
      <div className="absolute top-[30%] right-[10%] h-[300px] w-[300px] rounded-full bg-fuchsia-500/10 blur-[100px] z-0" />

      <header className="flex h-16 items-center px-4 lg:px-8 border-b border-white/20 bg-white/50 backdrop-blur-xl dark:bg-slate-900/50 dark:border-white/10 sticky top-0 z-50">
        <Link className="flex items-center justify-center gap-2" href="#">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 text-white shadow-lg shadow-indigo-500/30">
            <BookOpen className="h-5 w-5" />
          </div>
          <span className="font-bold text-xl tracking-tight text-foreground">AcademAide</span>
        </Link>
        <nav className="ml-auto flex gap-4 sm:gap-6">
          <Link className="text-sm font-medium hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors" href="#">
            Features
          </Link>
          <Link className="text-sm font-medium hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors" href="#">
            About
          </Link>
        </nav>
      </header>
      <main className="flex-1 z-10">
        <section className="w-full py-20 md:py-32 lg:py-40 flex flex-col items-center justify-center text-center px-4">
          <div className="container px-4 md:px-6 relative">
            <div className="flex flex-col items-center space-y-6">
              <div className="space-y-4 max-w-4xl">
                <h1 className="text-4xl font-extrabold tracking-tight sm:text-5xl md:text-6xl lg:text-7xl">
                  Your <span className="gradient-text">AI-Powered</span> <br className="hidden sm:inline" />
                  Academic Companion
                </h1>
                <p className="mx-auto max-w-[700px] text-muted-foreground md:text-xl leading-relaxed">
                  Personalized guidance, smart timetables, and instant academic support at your fingertips. Experience the future of learning today.
                </p>
              </div>
              <div className="flex flex-col sm:flex-row gap-4 min-w-[300px] justify-center pt-4">
                <Button asChild size="lg" className="rounded-full shadow-lg shadow-indigo-500/25 bg-indigo-600 hover:bg-indigo-700 h-12 px-8 text-base">
                  <Link href="/login">Get Started</Link>
                </Button>
                <Button asChild variant="outline" size="lg" className="rounded-full border-gray-200 bg-white/50 backdrop-blur-sm hover:bg-white hover:text-indigo-600 dark:border-white/10 dark:bg-slate-900/50 dark:hover:bg-slate-800 h-12 px-8 text-base">
                  <Link href="/signup">Create Account</Link>
                </Button>
              </div>
            </div>
          </div>
        </section>
      </main>
      <footer className="flex flex-col gap-2 sm:flex-row py-8 w-full shrink-0 items-center px-4 md:px-6 border-t border-gray-100 dark:border-white/5 bg-white/30 backdrop-blur-sm">
        <p className="text-xs text-muted-foreground">© 2024 AcademAide. All rights reserved.</p>
        <nav className="sm:ml-auto flex gap-4 sm:gap-6">
          <Link className="text-xs hover:text-indigo-600 underline-offset-4" href="#">
            Terms of Service
          </Link>
          <Link className="text-xs hover:text-indigo-600 underline-offset-4" href="#">
            Privacy
          </Link>
        </nav>
      </footer>
    </div>
  )
}
