import { Toaster } from "@/components/ui/toaster"
import '../App.css'
import wallpaper from '../login_page_wallpaper.jpg'
import { GitHubLogoIcon } from "@radix-ui/react-icons"
import { Button } from "@/components/ui/button"

export default function Login() {
    const handleGitHubLogin = () => {
        window.location.href = document.location.origin + '/auth/oauth'
    }

    const handleGitLabLogin = () => {
        // Implement GitLab OAuth login logic here
        window.location.href = document.location.origin + '/auth/oauth'
        console.log("GitLab login clicked")
    }

    const handleGoogleLogin = () => {
        // Implement Google OAuth login logic here
        window.location.href = document.location.origin + '/auth/oauth'
        console.log("Google login clicked")
    }

    return (
        <div className="w-full lg:grid lg:h-screen lg:grid-cols-2">
            <div className="flex items-center justify-center py-12">
                <div className="mx-auto grid w-[350px] gap-6">
                    <div className="grid gap-2 text-center">
                        <h1 className="text-3xl font-bold">Login</h1>
                        <p className="text-balance text-muted-foreground">
                            Choose a provider to login to your account
                        </p>
                    </div>
                    <div className="grid gap-4">
                        <Button variant="outline" className="w-full" onClick={handleGitHubLogin} disabled={true}>
                            <GitHubLogoIcon className="mr-2 h-4 w-4" />
                            Login with GitHub
                        </Button>
                        <Button variant="outline" className="w-full" onClick={handleGitLabLogin} disabled={true}>
                            <svg className="mr-2 h-4 w-4" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M23.955 13.587l-1.342-4.135-2.664-8.189a.455.455 0 00-.867 0L16.418 9.45H7.582L4.918 1.263a.455.455 0 00-.867 0L1.386 9.45.044 13.587a.924.924 0 00.331 1.03L12 23.054l11.625-8.436a.92.92 0 00.33-1.031" />
                            </svg>
                            Login with GitLab
                        </Button>
                        <Button variant="outline" className="w-full" onClick={handleGoogleLogin}>
                            <svg className="mr-2 h-4 w-4" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4" />
                                <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
                                <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05" />
                                <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
                                <path d="M1 1h22v22H1z" fill="none" />
                            </svg>
                            Login with Google
                        </Button>
                    </div>
                    <div className="mt-4 text-center text-sm">
                        Don&apos;t have an account?{" "}
                        <a href="#" className="underline">
                            Sign up
                        </a>
                    </div>
                </div>
                <Toaster />
            </div>
            <div className="hidden bg-muted lg:block bg-cover bg-center" style={{ backgroundImage: `url(${wallpaper})` }}></div>
        </div>
    )
}


