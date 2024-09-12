import { Toaster } from "@/components/ui/toaster"
import { Search, Copy } from "lucide-react"
import '../App.css'
import { PieChart, Pie, Cell, ResponsiveContainer } from 'recharts'

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Input } from "@/components/ui/input"
import { useToast } from "@/hooks/use-toast"
import { useState ,useEffect } from 'react'

export default function Dashbobard() {
  const [apiKey, setApiKey] = useState<string>("You dont have an API key.")
  const [pictureUrl, setPictureUrl] = useState<string>("")
  const [requestsSent, setRequestsSent] = useState(75) // Example value
  const { toast } = useToast()
  const [isDialogOpen, setIsDialogOpen] = useState(false)

  function getCookie(cname: string) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for(let i = 0; i <ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) == ' ') {
        c = c.substring(1);
      }
      if (c.indexOf(name) == 0) {
        return c.substring(name.length, c.length);
      }
    }
    return "";
  }

  useEffect(() => {
    const userInfo = getCookie("user").replace('"', '')
    if (userInfo == "") {
        return 
    }

    setRequestsSent(0)
    const cookieData = new Map<string, string>();
    userInfo.split(',').map(e => {
        const [key, value] = e.split('=')
        cookieData.set(key, value)

    })
    const apiKey = cookieData.get("apiKey") || ""
    const picture = cookieData.get("picture") || ""
    setApiKey(apiKey)
    setPictureUrl(picture)
  }, []);

  const copyToClipboard = () => {
    navigator.clipboard.writeText(apiKey).then(() => {
      toast({
        className: "text-white bg-green-700",
        title: "Copied!",
        description: "API key copied to clipboard.",
      })
    })
  }

  const handleKeyGeneration = () => {
    const fetchData = async () => {
        const response = await fetch(document.location.origin + "/renew")
        if (!response.ok) {
          throw new Error(`Response status: ${response.status}`);
        }

        const json = await response.json();
        setApiKey(json["api_key"] || "")
        toast({
          className: "text-white bg-green-700",
          title: "New key generated !",
          description: "You account have a new API Key.",
        })
    }
    fetchData().catch(console.error);
    setIsDialogOpen(false)
  }

  const handleLogout = () => {
        window.location.href = document.location.origin + "/logout"
    setIsDialogOpen(false)
  }
  const data = [
    { name: 'Requests Sent', value: requestsSent },
    { name: 'Remaining', value: 100 - requestsSent },
  ]

  const COLORS = ['#8884d8', '#f3f4f6']

  return (
    <div className="flex min-h-screen w-full flex-col">
      <header className="sticky top-0 flex h-16 items-center gap-4 border-b bg-background px-4 md:px-6">
        <div className="flex w-full items-center gap-4 md:ml-auto md:gap-2 lg:gap-4">
          <form className="ml-auto flex-1 sm:flex-initial">
            <div className="relative">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                type="search"
                placeholder="Search products..."
                className="pl-8 sm:w-[300px] md:w-[200px] lg:w-[300px]"
              />
            </div>
          </form>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="secondary" size="icon" className="rounded-full">
                <Avatar>
                    <AvatarImage src={pictureUrl} />
                    <AvatarFallback>CN</AvatarFallback>
                </Avatar>
                <span className="sr-only">Toggle user menu</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>My Account</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => handleLogout()}>Logout</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </header>
      <main className="flex min-h-[calc(100vh_-_theme(spacing.16))] flex-1 flex-col gap-4 bg-muted/40 p-4 md:gap-8 md:p-10">
        <div className="mx-auto grid w-full max-w-6xl gap-2">
          <h1 className="text-3xl font-semibold">Settings</h1>
        </div>
        <div className="mx-auto grid w-full max-w-6xl gap-6 md:grid-cols-[1fr_2fr]">
          <Card className="flex flex-col">
            <CardHeader>
              <CardTitle>Requests</CardTitle>
              <CardDescription>
                Number of requests sent in the last 30 days
              </CardDescription>
            </CardHeader>
            <CardContent className="flex-grow flex justify-center items-center">
              <div className="w-64 h-64 relative">
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={data}
                      cx="50%"
                      cy="50%"
                      innerRadius={60}
                      outerRadius={80}
                      fill="#8884d8"
                      paddingAngle={5}
                      dataKey="value"
                    >
                      {data.map((entry) => (
                        <Cell key={`cell-${entry.value}`} fill={COLORS[entry.value % COLORS.length]} />
                      ))}
                    </Pie>
                  </PieChart>
                </ResponsiveContainer>
                <div className="absolute inset-0 flex items-center justify-center">
                  <p className="text-2xl font-bold">{requestsSent}</p>
                </div>
              </div>
            </CardContent>
            <CardFooter className="justify-between">
              <div>{requestsSent} / 100 requests</div>
              <Button variant="outline">View Details</Button>
            </CardFooter>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle>API Key</CardTitle>
              <CardDescription>
                Generate an key to interact with the API
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex items-center space-x-2">
                <Input
                  value={apiKey}
                  readOnly
                  className="flex-grow text-gray-500"
                />
                <Button
                  variant="outline"
                  size="icon"
                  onClick={copyToClipboard}
                  className="p-0"
                >
                  <Copy className="left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                </Button>
              </div>
            </CardContent>
            <CardFooter className="mt-auto border-t px-6 py-4">
              <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
                <DialogTrigger asChild>
                  <Button>Generate API Key</Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Confirm Update</DialogTitle>
                    <DialogDescription>
                      Are you sure you want to generate a new API key ? This action cannot be undone.
                    </DialogDescription>
                  </DialogHeader>
                  <DialogFooter>
                    <Button variant="outline" onClick={() => setIsDialogOpen(false)}>
                      Cancel
                    </Button>
                    <Button onClick={handleKeyGeneration}>Confirm</Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
            </CardFooter>
          </Card>
        </div>
      <Toaster />
      </main>
    </div>
  )
}
