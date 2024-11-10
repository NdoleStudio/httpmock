import { SignIn } from '@clerk/nextjs'
import {Box} from "@primer/react";

export default function Page() {
    return (
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: 'calc(100vh - 200px)'}}>
            <SignIn />
        </Box>
    )
}