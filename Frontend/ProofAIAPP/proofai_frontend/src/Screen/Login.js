import React from "react";
import { useNavigate } from "react-router-dom";
import LoginBox from "../Components/LoginBox";
import Footer from "../Components/Footer";
import SignupBox from "../Components/SignupBox";

const Login = () => {

    const navigate = useNavigate()
    const handleLogin = () => {
        navigate('/Home')
    }

    const [isLogin, setIsLogin] = React.useState(true)
    const handleLoginSignup = () => {
        setIsLogin(!isLogin)
    }

    return (
        <div>
            {isLogin && <div style={styles.LoginBox}> <LoginBox onPress={handleLogin} handleLoginSignup={handleLoginSignup} /> </div>}
            {!isLogin && <div style={styles.LoginBox}> <SignupBox handleLoginSignup={handleLoginSignup} onPress={handleLogin} /> </div>}
        </div>
    );
};

export default Login;

const styles = ({

    LoginBox: {
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        width: '100%',
        height: '100vh',
        minHeight: '100vh',
        backgroundColor: '#f4f6f7',
        boxSizing: 'border-box',
        paddingBottom: '100px',
    },
})
