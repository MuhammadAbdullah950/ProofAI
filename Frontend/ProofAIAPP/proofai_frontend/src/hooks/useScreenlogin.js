import React from "react";
import { useNavigate } from "react-router-dom";
import LoginBox from "../Components/LoginBox";
import Footer from "../Components/Footer";
import SignupBox from "../Components/SignupBox";

const useScreenlogin = () => {

    const navigate = useNavigate()
    const handleLogin = () => {
        navigate('/Home')
    }

    const [isLogin, setIsLogin] = React.useState(true)
    const handleLoginSignup = () => {
        setIsLogin(!isLogin)
    }
    const [boxHeight, setBoxHeight] = useState("auto");
    useEffect(() => {
        const updateHeight = () => {
            const viewportHeight = window.innerHeight;
            const minHeight = 400; // Minimum height for the auth box
            const padding = 48; // 24px top + 24px bottom
            const calculatedHeight = viewportHeight - padding;

            setBoxHeight(Math.max(minHeight, calculatedHeight));
        };

        updateHeight();
        window.addEventListener('resize', updateHeight);

        return () => window.removeEventListener('resize', updateHeight);
    }, []);


    return {
        isLogin,
        handleLogin,
        handleLoginSignup,
        boxHeight,
    }
};

export default useScreenlogin;
